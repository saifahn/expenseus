import { PropsWithChildren } from 'react';

export default function SharedLayout({ children }: PropsWithChildren<{}>) {
  return (
    <>
      <section className="mt-4">{children}</section>
    </>
  );
}
