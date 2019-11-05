export interface StorageFile {
  name: string;
  size: number;
  isDirectory: boolean;
  updatedAt: Date | string;
}
