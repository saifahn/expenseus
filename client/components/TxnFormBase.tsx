import { enUSCategories } from 'data/categories';
import React from 'react';
import { UseFormRegisterReturn } from 'react-hook-form';

type InputProps = {
  title?: string;
  txnNameInputProps: UseFormRegisterReturn;
  amountInputProps: UseFormRegisterReturn;
  dateInputProps: UseFormRegisterReturn;
  categoryInputProps: UseFormRegisterReturn;
  children?: React.ReactNode;
  onSubmit: () => void;
};

export default function TxnFormBase({
  title,
  txnNameInputProps,
  amountInputProps,
  dateInputProps,
  categoryInputProps,
  children,
  onSubmit,
}: InputProps) {
  return (
    <form
      className="border-4 p-6"
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit();
      }}
    >
      {title && <h3 className="text-lg font-semibold">{title}</h3>}
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="name">
          Name
        </label>
        <input
          {...txnNameInputProps}
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
          id="name"
          type="text"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="amount">
          Amount
        </label>
        <input
          {...amountInputProps}
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
          id="amount"
          type="text"
          inputMode="numeric"
          pattern="[0-9]*"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="date">
          Date
        </label>
        <input
          {...dateInputProps}
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
          type="date"
          id="date"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold">Category</label>
        <select
          {...categoryInputProps}
          className="mt-2 block rounded bg-white bg-clip-padding bg-no-repeat px-3 py-2 text-base font-normal text-gray-700 outline outline-1 transition ease-in-out focus:border-indigo-600 focus:bg-white focus:text-gray-700"
        >
          {enUSCategories.map((category) => (
            <option key={category.key} value={category.key}>
              {category.value}
            </option>
          ))}
        </select>
      </div>
      {children}
    </form>
  );
}
