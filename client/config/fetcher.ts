export const fetcher = (url) =>
  fetch(url, { credentials: 'include' }).then((res) => res.json());
