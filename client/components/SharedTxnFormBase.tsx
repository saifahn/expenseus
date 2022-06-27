import { enUSCategories } from 'data/categories';
import { Tracker } from 'pages/shared/trackers';
import React from 'react';
import { UseFormRegisterReturn } from 'react-hook-form';

type Props = {
  title?: string;
  shopInputProps: UseFormRegisterReturn;
  amountInputProps: UseFormRegisterReturn;
  dateInputProps: UseFormRegisterReturn;
  categoryInputProps: UseFormRegisterReturn;
  payerInputProps: UseFormRegisterReturn;
  settledInputProps: UseFormRegisterReturn;
  tracker: Tracker;
  children?: React.ReactNode;
  onSubmit: () => void;
};

export default function SharedTxnFormBase({
  title,
  shopInputProps,
  amountInputProps,
  dateInputProps,
  categoryInputProps,
  payerInputProps,
  settledInputProps,
  tracker,
  children,
  onSubmit,
}: Props) {
  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit();
      }}
      className="border-4 p-6"
    >
      {title && <h3 className="text-lg font-semibold">{title}</h3>}
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="shop">
          Shop
        </label>
        <input
          {...shopInputProps}
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
          type="text"
          id="shop"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="amount">
          Amount
        </label>
        <input
          {...amountInputProps}
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
          type="text"
          inputMode="numeric"
          pattern="[0-9]*"
          id="amount"
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
        <label className="block font-semibold" htmlFor="payer">
          Payer
        </label>
        <select
          {...payerInputProps}
          className="mt-2 block rounded bg-white bg-clip-padding bg-no-repeat px-3 py-2 text-base font-normal text-gray-700 outline outline-1 transition ease-in-out focus:border-indigo-600 focus:bg-white focus:text-gray-700"
        >
          {tracker.users.map((userId) => (
            <option key={userId} value={userId}>
              {userId}
            </option>
          ))}
        </select>
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="settled">
          Settled?
        </label>
        <input {...settledInputProps} type="checkbox" id="settled" />
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
