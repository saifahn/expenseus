export async function withAsyncTryCatch<T>(
  promise: Promise<T>,
): Promise<[T | null, unknown]> {
  try {
    const data = await promise;
    return [data, null];
  } catch (err) {
    return [null, err];
  }
}

export function withTryCatch<T extends (...args: any) => any>(
  fn: T,
): [ReturnType<T> | null, unknown] {
  try {
    const data = fn();
    return [data, null];
  } catch (err) {
    return [null, err];
  }
}
