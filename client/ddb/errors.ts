export class ItemDoesNotExistError extends Error {
  constructor(message?: string) {
    super(message);
  }
}
