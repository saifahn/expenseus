import { CategoryKey, enUSCategories } from 'data/categories';
import React from 'react';
import { UseFormRegisterReturn } from 'react-hook-form';

export type TxnFormInputs = {
  location: string;
  amount: number;
  date: string;
  category: CategoryKey;
  details: string;
};

export function createTxnFormData(data: TxnFormInputs) {
  const formData = new FormData();
  formData.append('location', data.location);
  formData.append('details', data.details);
  formData.append('amount', data.amount.toString());
  formData.append('category', data.category);

  const unixDate = new Date(data.date).getTime();
  formData.append('date', unixDate.toString());
  return formData;
}

type Props = {
  title?: string;
  locationInputProps: UseFormRegisterReturn;
  amountInputProps: UseFormRegisterReturn;
  dateInputProps: UseFormRegisterReturn;
  categoryInputProps: UseFormRegisterReturn;
  detailsInputProps: UseFormRegisterReturn;
  children?: React.ReactNode;
  onSubmit: () => void;
};

export default function TxnFormBase({
  title,
  locationInputProps,
  detailsInputProps,
  amountInputProps,
  dateInputProps,
  categoryInputProps,
  children,
  onSubmit,
}: Props) {
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
        <label className="block font-semibold" htmlFor="location">
          Location
        </label>
        <input
          {...locationInputProps}
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
          id="location"
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
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="details">
          Details
        </label>
        <input
          {...detailsInputProps}
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
          id="details"
          type="text"
        />
      </div>
      {children}
    </form>
  );
}
