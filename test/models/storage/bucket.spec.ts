import { expect } from "chai";
import fs from "fs";
import S3 from "aws-sdk/clients/s3";
import sinon from "sinon";

import { StorageBucket } from "../../../src/models/storage/bucket";
import S3Client from "../../../src/libs/s3";

describe("StorageBucket", () => {

  const bucket = new StorageBucket();
  bucket.name = "mybucket";
  bucket.isPublic = true;
  bucket.createdAt = new Date();

  context("#readFile", () => {
    const output: S3.GetObjectOutput = { Body: Buffer.from("example", "utf-8") };
    const spy = sinon.stub(S3Client.prototype, "read").returns(new Promise((resolve) => resolve(output)));

    beforeEach(() => {
      const keyHash = Buffer.from("mykey").toString("base64");
      const cacheFile = `${process.env.LUPPITER_STORAGE_CACHE_PATH || "/tmp"}/mybucket/${keyHash}`;
      if (fs.existsSync(cacheFile)) {
        fs.unlinkSync(cacheFile);
      }

      spy.resetHistory();
    });

    it("first try should download from s3", async () => {
      expect((await bucket.readFile("mykey")).toString()).to.equal(output.Body.toString());
      expect(spy.calledOnceWith("mybucket/mykey")).to.be.true;
    });

    it("second try should use cache", async () => {
      expect((await bucket.readFile("mykey")).toString()).to.equal(output.Body.toString());
      expect((await bucket.readFile("mykey")).toString()).to.equal(output.Body.toString());
      expect(spy.calledOnceWith("mybucket/mykey")).to.be.true;
    });
  });

  context("#writeFile", () => {
    const output: S3.PutObjectOutput = {};
    const spy = sinon.stub(S3Client.prototype, "write").returns(new Promise((resolve) => resolve(output)));
    after(() => spy.restore());

    const body = Buffer.from("example", "utf-8");

    it("should success", async () => {
      expect(await bucket.writeFile("mykey", body)).not.to.throw;
      expect(spy.calledOnceWith("mykey", body)).to.be.true;
    });
  });

  context("#listFiles", () => {
    const output: S3.ListObjectsV2Output = {Contents: [], CommonPrefixes: []};
    const spy = sinon.stub(S3Client.prototype, "list").returns(new Promise((resolve) => resolve(output)));

    it("should success", async () => {
      expect(await bucket.listFiles()).not.to.throw;
    });

    spy.restore();
  });

});
