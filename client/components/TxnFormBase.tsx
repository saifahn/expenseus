import {
  categories,
  SubcategoryKey,
  mainCategories,
  mainCategoryKeys,
  categoryNameFromKeyEN,
} from 'data/categories';
import { CreateTxnPayload } from 'pages/api/v1/transactions';
import React from 'react';
import { UseFormRegister } from 'react-hook-form';
import { plainDateStringToEpochSec } from 'utils/dates';

export type TxnFormInputs = {
  location: string;
  amount: number;
  date: string;
  category: SubcategoryKey;
  details: string;
};

export function makeCreateTxnPayload(
  data: TxnFormInputs,
): Omit<CreateTxnPayload, 'userId'> {
  return {
    location: data.location,
    details: data.details,
    amount: Number(data.amount),
    category: data.category,
    date: plainDateStringToEpochSec(data.date),
  };
}

export function createTxnFormData(data: TxnFormInputs) {
  const formData = new FormData();
  formData.append('location', data.location);
  formData.append('details', data.details);
  formData.append('amount', data.amount.toString());
  formData.append('category', data.category);

  const unixSeconds = plainDateStringToEpochSec(data.date);
  formData.append('date', unixSeconds.toString());
  return formData;
}

type Props = {
  title?: string;
  register: UseFormRegister<TxnFormInputs>;
  onSubmit: () => void;
  children?: React.ReactNode;
};

export default function TxnFormBase({
  title,
  register,
  onSubmit,
  children,
}: Props) {
  return (
    <form
      className="bg-white py-3"
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit();
      }}
    >
      {title && <h3 className="text-lg font-bold lowercase">{title}</h3>}
      <div className="mt-3">
        <label
          className="block font-semibold lowercase text-slate-600"
          htmlFor="amount"
        >
          Amount
        </label>
        <input
          {...register('amount', {
            min: { value: 1, message: 'Please input a positive amount' },
            required: 'Please input an amount',
          })}
          className="focus:border-violet block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 text-center text-xl placeholder-slate-400 focus:ring-0"
          id="amount"
          type="text"
          inputMode="numeric"
          pattern="[0-9]*"
          placeholder="93872円"
        />
      </div>
      <div className="mt-5">
        <label
          className="block font-semibold lowercase text-slate-600"
          htmlFor="location"
        >
          Description
        </label>
        <input
          {...register('location', {
            required: 'Please input a description',
          })}
          className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 placeholder-slate-400 focus:ring-0"
          id="location"
          type="text"
          placeholder="what and where?"
        />
      </div>
      <div className="mt-5">
        <label
          className="block font-semibold lowercase text-slate-600"
          htmlFor="date"
        >
          Date
        </label>
        <input
          {...register('date', { required: 'Please input a date' })}
          className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 placeholder-slate-400 focus:ring-0"
          type="date"
          id="date"
        />
      </div>
      <div className="mt-5">
        <label className="block font-semibold lowercase text-slate-600">
          Category
        </label>
        <select
          {...register('category')}
          className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 lowercase placeholder-slate-400 focus:ring-0"
        >
          {mainCategoryKeys.map((mainKey) => (
            <optgroup
              key={mainKey}
              label={`${mainCategories[mainKey].emoji} ${mainCategories[mainKey].en_US}`}
            >
              {categories[mainKey].map((subKey) => (
                <option key={subKey} value={subKey}>
                  {categoryNameFromKeyEN(subKey)}
                </option>
              ))}
            </optgroup>
          ))}
        </select>
      </div>
      <div className="mt-5">
        <label
          className="block font-semibold lowercase text-slate-600"
          htmlFor="details"
        >
          Details
        </label>
        <input
          {...register('details')}
          className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 placeholder-slate-400 focus:ring-0"
          id="details"
          type="text"
          placeholder="any other details?"
        />
      </div>
      {children}
    </form>
  );
}
