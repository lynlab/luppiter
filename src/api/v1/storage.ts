import { Request, Response } from "express";
import expressContext from "express-http-context";

import { ApiKey } from "../../models/auth/api_key";
import { StorageBucket } from "../../models/storage/bucket";

// GET /vulcan/storage/buckets
//
// Required permission: `Storage::Write`
async function listBuckets(_: Request, res: Response) {
  const apiKey: ApiKey = expressContext.get("request:api_key");

  const buckets = await StorageBucket.find({ member: apiKey.member });
  res.json(buckets.map((bucket) => bucket.toJson()));
}

// POST /vulcan/storage/buckets
//
// Required permission: `Storage::Write`
// Request Body:
// {
//   "name": "string",
//   "isPublic": "boolean?"
// }
// Errors:
//   - duplicated_entry(400) : given bucket name has already exists
async function createBucket(req: Request, res: Response) {
  if (await StorageBucket.findOne({ name: req.body.name })) {
    res.status(400).json({ error: "duplicated_entry" });
    return;
  }

  const apiKey: ApiKey = expressContext.get("request:api_key");
  const bucket = new StorageBucket();
  bucket.member = apiKey.member;
  bucket.name = req.body.name;
  bucket.isPublic = req.body.isPublic === true;
  await bucket.save();

  res.json(bucket.toJson());
}

// PUT /vulcan/storage/buckets/:name
//
// Required Permission: `Storage::Write`
// Request Body:
// {
//   "isPublic": "boolean"
// }
async function updateBucket(req: Request, res: Response) {
  const apiKey: ApiKey = expressContext.get("request:api_key");
  const bucket = await StorageBucket.findOne({ where: { name: req.params.name }, relations: ["member"] });
  if (!bucket || bucket.member.id !== apiKey.member.id) {
    res.sendStatus(401);
    return;
  }

  bucket.isPublic = req.body.isPublic === true;
  await bucket.save();

  res.json(bucket.toJson());
}

// DELETE /vulcan/storage/buckets/:name
//
// Required Permission: `Storage::Wrtie`
async function deleteBucket(req: Request, res: Response) {
  const apiKey: ApiKey = expressContext.get("request:api_key");
  const bucket = await StorageBucket.findOne({ where: { name: req.params.name }, relations: ["member"] });
  if (!bucket || bucket.member.id !== apiKey.member.id) {
    res.sendStatus(401);
    return;
  }

  await bucket.remove();

  res.json(bucket.toJson());
}

// GET /vulcan/storage/buckets/:name/files
//
// Required permission: `Storage::Read`
// Query params:
//   - prefix(string?) : limits name of keys with the specified prefix
//   - cursor(string?) : base64 + url encoded file name
async function listBucketFiles(req: Request, res: Response) {
  const apiKey: ApiKey = expressContext.get("request:api_key");
  const bucket = await StorageBucket.findOne({ where: { name: req.params.name }, relations: ["member"] });
  if (!bucket || bucket.member.id !== apiKey.member.id) {
    res.sendStatus(401);
    return;
  }

  const { prefix, cursor } = req.query;
  let decodedCursor: string;
  if (cursor) {
    decodedCursor = Buffer.from(cursor, "base64").toString("utf-8");
  }

  res.json(await bucket.listFiles(prefix, decodedCursor));
}

export default {
  listBuckets,
  createBucket,
  updateBucket,
  deleteBucket,
  listBucketFiles,
};
