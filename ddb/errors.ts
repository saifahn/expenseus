export class ItemDoesNotExistError extends Error {
  constructor(message?: string) {
    super(message);
  }
}

export class UserAlreadyExistsError extends Error {
  constructor() {
    super();
  }
}
