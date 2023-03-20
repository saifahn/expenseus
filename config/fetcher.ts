export class HttpException extends Error {
  code: number;
  info: any;

  constructor({
    message,
    code,
    info,
  }: {
    message: string;
    code: number;
    info: any;
  }) {
    super(message);
    this.code = code;
    this.info = info;
  }
}

export async function fetcher(url: string) {
  const res = await fetch(url, { credentials: 'include' });

  if (!res.ok) {
    let info = await res.text();
    const error = new HttpException({
      message: 'An error occurred while fetching the data.',
      code: res.status,
      info,
    });
    throw error;
  }

  return res.json();
}
