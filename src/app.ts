import cors from "cors";
import express, { RequestHandler } from "express";
import expressFileupload from "express-fileupload";
import expressContext from "express-http-context";

import storage from "./api/storage";
import vulcanAuth from "./api/v1/auth";
import vulcanCerts from "./api/v1/certs";
import vulcanCloudContainer from "./api/v1/cloudcontainer";
import vulcanHosting from "./api/v1/hosting";
import vulcanStorage from "./api/v1/storage";
import { ApiKey } from "./models/auth/api_key";

const app = express();

app.use(express.json());
app.use(expressContext.middleware);
app.use(expressFileupload({
  limits: { fileSize: 10 * 1024 * 1024 },
}));

const corsOptions = {
  origin: (process.env.LUPPITER_ALLOWED_ORIGINS || "http://127.0.0.1:8080").split(","),
};
app.use(cors(corsOptions));

function permitted(handle: RequestHandler, permission: string, optional: boolean = false): RequestHandler {
  return async (req, res, next) => {
    const key = req.header("X-Api-Key");
    if (!key) {
      if (optional) {
        return handle(req, res, next);
      }

      res.sendStatus(401);
      return;
    }

    const apiKey = await ApiKey.findOne({ where: { key }, relations: ["member", "permissions"] });
    if (!apiKey || !apiKey.hasPermission(permission)) {
      res.sendStatus(401);
      return;
    }

    expressContext.set("request:api_key", apiKey);
    return handle(req, res, next);
  };
}

// V1 (Vulcan)
app.get("/vulcan/auth/me", vulcanAuth.getMe);
app.get("/vulcan/auth/api_keys", vulcanAuth.listApiKeys);
app.post("/vulcan/auth/api_keys", vulcanAuth.createApiKey);
app.delete("/vulcan/auth/api_keys/:key", vulcanAuth.deleteApiKey);
app.get("/vulcan/auth/api_keys/:key/permissions", vulcanAuth.listKeyPermissions);
app.post("/vulcan/auth/api_keys/:key/permissions", vulcanAuth.addKeyPermission);
app.delete("/vulcan/auth/api_keys/:key/permissions", vulcanAuth.removeKeyPermission);
app.get("/vulcan/auth/permissions", vulcanAuth.listPermissions);

app.get("/vulcan/certs/certificates", permitted(vulcanCerts.listCertificates, "Certs::Read"));
app.post("/vulcan/certs/certificates", permitted(vulcanCerts.createCertificate, "Certs::Write"));

app.get("/vulcan/cloudcontainer/tasks", permitted(vulcanCloudContainer.listTasks, "CloudContainer::*"));
app.post("/vulcan/cloudcontainer/tasks", permitted(vulcanCloudContainer.createTask, "CloudContainer::Write"));
app.put("/vulcan/cloudcontainer/tasks/:uuid", permitted(vulcanCloudContainer.updateTask, "CloudContainer::Write"));
app.delete("/vulcan/cloudcontainer/tasks/:uuid", permitted(vulcanCloudContainer.deleteTask, "CloudContainer::Write"));
app.post("/vulcan/cloudcontainer/tasks/:uuid/run", permitted(vulcanCloudContainer.runTask, "CloudContainer::Write"));

app.get("/vulcan/hosting/instances", permitted(vulcanHosting.listInstances, "Hosting::Read"));
app.post("/vulcan/hosting/instances", permitted(vulcanHosting.createInstance, "Hosting::Write"));
app.delete("/vulcan/hosting/instances/:uuid", permitted(vulcanHosting.deleteInstance, "Hosting::Write"));

app.get("/vulcan/storage/buckets", permitted(vulcanStorage.listBuckets, "Storage::Read"));
app.post("/vulcan/storage/buckets", permitted(vulcanStorage.createBucket, "Storage::Write"));
app.put("/vulcan/storage/buckets/:name", permitted(vulcanStorage.updateBucket, "Storage::Write"));
app.delete("/vulcan/storage/buckets/:name", permitted(vulcanStorage.deleteBucket, "Storage::Write"));
app.get("/vulcan/storage/buckets/:name/files", permitted(vulcanStorage.listBucketFiles, "Storage::Read"));

// File Storage
app.get("/storage/:namespace/:key(*)", permitted(storage.readFile, "Storage::Read", true));
app.post("/storage/:namespace/:key(*)", permitted(storage.writeFile, "Storage::Write", true));

app.get("/ping", (req, res) => res.send("pong"));

export default app;
