import fs from "fs";
import {
  BaseEntity, Column, CreateDateColumn, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn, UpdateDateColumn,
} from "typeorm";

import S3Client from "../../libs/s3";
import { Member } from "../auth/member";
import { StorageFile } from "./file";

@Entity({ name: "storage_buckets" })
export class StorageBucket extends BaseEntity {

  @PrimaryGeneratedColumn()
  public id: number;

  @ManyToOne((type) => Member)
  @JoinColumn({ name: "member_id" })
  public member: Member;

  @Column({ name: "name" })
  public name: string;

  @Column({ name: "is_public" })
  public isPublic: boolean;

  @CreateDateColumn({ name: "created_at" })
  public createdAt: Date;

  @UpdateDateColumn({ name: "updated_at" })
  public updatedAt: Date;

  private cachePath = process.env.LUPPITER_STORAGE_CACHE_PATH || "/tmp";
  private s3 = new S3Client("luppiter.lynlab.co.kr");

  public toJson(): object {
    return {
      name: this.name,
      isPublic: this.isPublic,
      createdAt: this.createdAt.toISOString(),
    };
  }

  /**
   * Get a file from bucket.
   * If cached file exists, return it immediately. If not, download the file from S3 and cache it.
   * @param key Key of file
   */
  public async readFile(key: string): Promise<Buffer> {
    const keyHash = Buffer.from(key).toString("base64");
    const cacheFile = `${this.cachePath}/${this.name}/${keyHash}`;
    if (fs.existsSync(cacheFile)) {
      return fs.readFileSync(cacheFile);
    }

    let fileBody: Buffer;
    try {
      const result = await this.s3.read(`${this.name}/${key}`);
      fileBody = result.Body as Buffer;
    } catch (e) {
      if (e.code === "NoSuchKey") {
        return null;
      }
      throw e;
    }

    const cachePath = `${this.cachePath}/${this.name}`;
    if (!fs.existsSync(cachePath)) {
      fs.mkdirSync(cachePath, { recursive: true, mode: 0o700 });
    }
    fs.writeFileSync(cacheFile, fileBody, { mode: 0o600 });

    return fileBody;
  }

  public async listFiles(prefix?: string, cursor?: string): Promise<StorageFile[]> {
    const bucketPrefix = prefix ? `${this.name}/${prefix}/` : `${this.name}/`;
    const bucketCursor = cursor ? `${this.name}/${cursor}` : null;
    const result = await this.s3.list(bucketPrefix, bucketCursor);

    const files = result.Contents.map((c) => {
      return {
        name: c.Key.replace(`${this.name}/`, ""),
        size: c.Size,
        isDirectory: false,
        updatedAt: c.LastModified.toISOString(),
      };
    }).filter((f) => f.size !== 0);
    const directories = result.CommonPrefixes.map((p) => {
      return { name: p.Prefix.replace(`${this.name}/`, ""), size: 0, isDirectory: true, updatedAt: null };
    });

    return files.concat(directories);
  }

  public async writeFile(key: string, body: Buffer) {
    this.s3.write(key, body);
  }

}
