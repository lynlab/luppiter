import { Request, Response } from "express";
import { UploadedFile } from "express-fileupload";
import expressContext from "express-http-context";
import fileType from "file-type";
import mimeTypes from "mime-types";

import { ApiKey } from "../models/auth/api_key";
import { StorageBucket } from "../models/storage/bucket";

async function readFile(req: Request, res: Response) {
  const { namespace, key } = req.params;
  if (!key) {
    res.sendStatus(404);
    return;
  }

  const bucket = await StorageBucket.findOne({ where: { name: namespace }, relations: ["member"] });
  const apiKey: ApiKey = expressContext.get("request:api_key");
  if (!bucket || (!bucket.isPublic && (!apiKey || bucket.member.id !== apiKey.member.id))) {
    res.sendStatus(401);
    return;
  }

  const fileBody = await bucket.readFile(key);
  if (!fileBody) {
    res.sendStatus(404);
    return;
  }

  const type = fileType(fileBody);
  let contentType: string;
  if (type) {
    contentType = type.mime;
  } else {
    contentType = mimeTypes.lookup(key) || "application/text";
  }

  res.header("Content-Type", contentType).send(fileBody);
}

async function writeFile(req: Request, res: Response) {
  const { namespace, key } = req.params;

  const bucket = await StorageBucket.findOne({ where: { name: namespace }, relations: ["member"] });
  const apiKey: ApiKey = expressContext.get("request:api_key");
  if (!bucket || !apiKey || bucket.member.id !== apiKey.member.id) {
    res.sendStatus(401);
    return;
  }

  const file = req.files.file as UploadedFile;
  try {
    await bucket.writeFile(key, file.data);
    res.sendStatus(201);
  } catch (e) {
    res.sendStatus(500);
  }
}

export default { readFile, writeFile };
