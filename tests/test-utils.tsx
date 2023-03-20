import { render, RenderOptions } from '@testing-library/react';
import { PropsWithChildren, ReactElement } from 'react';

const AllProviders = ({ children }: PropsWithChildren<{}>) => {
  return <>{{ children }}</>;
};

// override render method
const customRender = (
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
) => render(ui, { wrapper: AllProviders, ...options });

export { default as userEvent } from '@testing-library/user-event';
export * from '@testing-library/react';
export { customRender as render };
