import S3 from "aws-sdk/clients/s3";
import { expect } from "chai";
import faker from "faker";
import { Factory } from "rosie";
import sinon from "sinon";

import mockRequest from "../mock-request";
import S3Client from "../../../src/libs/s3";
import { Member } from "../../../src/models/auth/member";
import { ApiKey } from "../../../src/models/auth/api_key";
import { Permission } from "../../../src/models/auth/permission";
import { StorageBucket } from "../../../src/models/storage/bucket";

describe("vulcan storage apis", () => {
  let member: Member;
  let apiKey: ApiKey;
  let stranger: Member;
  let strangerKey: ApiKey;
  
  before(async () => {
    member = await Factory.build<Member>("Member");
    apiKey = await Factory.build<ApiKey>("ApiKey", {
      member,
      permissions: [await Permission.findOne({ key: "Storage::*" })],
    });

    stranger = await Factory.build<Member>("Member");
    strangerKey = await Factory.build<ApiKey>("ApiKey", {
      member: stranger,
      permissions: [await Permission.findOne({ key: "Storage::*" })],
    });
  });

  describe("GET /vulcan/storage/buckets", () => {
    let bucket: StorageBucket;
    before(async () => {
      bucket = await Factory.build<StorageBucket>("StorageBucket", { member });
    });

    it("success", (done) => {
      mockRequest.get('/vulcan/storage/buckets')
        .set("X-Api-Key", apiKey.key)
        .then((res) => {
          expect(res.status).to.equal(200);
          expect(res.body.length).to.equal(1);
          expect(res.body[0].name).to.equal(bucket.name);
          expect(res.body[0].isPublic).to.equal(bucket.isPublic);
          expect(res.body[0].createdAt).to.equal(bucket.createdAt.toISOString());
          done();
        })
        .catch(e => done(e));
    });
  });

  describe("POST /vulcan/storage/buckets", () => {
    const expectedName = faker.random.uuid();

    it("success", (done) => {
      mockRequest.post('/vulcan/storage/buckets')
        .set("X-Api-Key", apiKey.key)
        .send({ name: expectedName, isPublic: true })
        .then(async (res) => {
          expect(res.status).to.equal(200);
          expect(res.body.name).to.be.equal(expectedName);
          expect(res.body.isPublic).to.be.true;
          
          const bucket = await StorageBucket.findOne({ name: expectedName });
          expect(bucket.isPublic).to.be.true;

          done();
        })
        .catch(e => done(e));
    });

    it("failed for duplicated entry", (done) => {
      mockRequest.post('/vulcan/storage/buckets')
        .set("X-Api-Key", apiKey.key)
        .send({ name: expectedName, isPublic: true })
        .then(async (res) => {
          expect(res.status).to.equal(400);
          done();
        })
        .catch(e => done(e));
    });
  });

  describe("PUT /vulcan/storage/buckets/:name", () => {
    let bucket: StorageBucket;
    beforeEach(async () => {
      bucket = await Factory.build<StorageBucket>("StorageBucket", { member, isPublic: false });
    });

    it("success", (done) => {
      mockRequest.put(`/vulcan/storage/buckets/${bucket.name}`)
        .set("X-Api-Key", apiKey.key)
        .send({ isPublic: true })
        .then(async (res) => {
          expect(res.status).to.equal(200);
          expect(res.body.name).to.equal(bucket.name);
          expect(res.body.isPublic).to.be.true;

          bucket = await StorageBucket.findOne({ id: bucket.id });
          expect(bucket.isPublic).to.be.true;

          done();
        })
        .catch(e => done(e));
    });

    it("raise 401 for non-existing bucket", (done) => {
      mockRequest.put(`/vulcan/storage/buckets/${faker.random.uuid()}`)
        .set("X-Api-Key", apiKey.key)
        .then(res => {
          expect(res.status).to.equal(401);
          done();
        })
        .catch(e => done(e));
    });

    it("raise 401 for non-owner", (done) => {
      mockRequest.put(`/vulcan/storage/buckets/${bucket.name}`)
        .set("X-Api-Key", strangerKey.key)
        .then(res => {
          expect(res.status).to.equal(401);
          done();
        })
        .catch(e => done(e));
    });
  });

  describe("DELETE /vulcan/storage/buckets/:name", () => {
    let bucket: StorageBucket;
    beforeEach(async () => {
      bucket = await Factory.build<StorageBucket>("StorageBucket", { member });
    });

    it("success", (done) => {
      mockRequest.delete(`/vulcan/storage/buckets/${bucket.name}`)
        .set("X-Api-Key", apiKey.key)
        .then(async (res) => {
          expect(res.status).to.equal(200);
          expect(res.body.name).to.equal(bucket.name);

          bucket = await StorageBucket.findOne({ name: bucket.name });
          expect(bucket).to.be.undefined;
          done();
        })
        .catch(e => done(e));
    });

    it("raise 401 for non-existing bucket", (done) => {
      mockRequest.delete(`/vulcan/storage/buckets/${faker.random.uuid()}`)
        .set("X-Api-Key", apiKey.key)
        .then(res => {
          expect(res.status).to.equal(401);
          done();
        })
        .catch(e => done(e));
    });

    it("raise 401 for non-owner", (done) => {
      mockRequest.delete(`/vulcan/storage/buckets/${bucket.name}`)
        .set("X-Api-Key", strangerKey.key)
        .then(res => {
          expect(res.status).to.equal(401);
          done();
        })
        .catch(e => done(e));
    });
  });

  describe("GET /vulcan/storage/buckets/:name/files", () => {
    const output: S3.ListObjectsV2Output = { Contents: [], CommonPrefixes: [] };
    const spy = sinon.stub(S3Client.prototype, "list").returns(new Promise((resolve) => resolve(output)));

    let bucket: StorageBucket;
    beforeEach(async () => {
      bucket = await Factory.build<StorageBucket>("StorageBucket", { member });
    });

    it("success", (done) => {
      mockRequest.get(`/vulcan/storage/buckets/${bucket.name}/files`)
        .set("X-Api-Key", apiKey.key)
        .then(res => {
          expect(res.status).to.equal(200);
          done();
        })
        .catch(e => done(e));
    });

    spy.restore();
  });
});
