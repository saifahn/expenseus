import {
  categories,
  SubcategoryKey,
  mainCategories,
  mainCategoryKeys,
  subcategories,
  categoryNameFromKeyEN,
} from 'data/categories';
import { Tracker } from 'pages/shared/trackers';
import React, { useState } from 'react';
import { UseFormRegister } from 'react-hook-form';
import { plainDateStringToEpochSec } from 'utils/dates';

export type SharedTxnFormInputs = {
  location: string;
  amount: number;
  date: string;
  participants: string;
  payer: string;
  details: string;
  category: SubcategoryKey;
  settled?: boolean;
  split?: string;
};

export function createSharedTxnFormData(data: SharedTxnFormInputs) {
  const formData = new FormData();

  formData.append('location', data.location);
  formData.append('amount', data.amount.toString());
  if (!data.settled) formData.append('unsettled', 'true');
  formData.append('category', data.category);
  formData.append('payer', data.payer);
  formData.append('details', data.details);

  const unixDate = plainDateStringToEpochSec(data.date);
  formData.append('date', unixDate.toString());
  formData.append('split', data.split);
  return formData;
}

type Props = {
  title?: string;
  tracker: Tracker;
  register: UseFormRegister<SharedTxnFormInputs>;
  onSubmit: () => void;
  children?: React.ReactNode;
};

export default function SharedTxnFormBase({
  title,
  tracker,
  register,
  onSubmit,
  children,
}: Props) {
  const [hasCustomSplit, setHasCustomSplit] = useState(false);

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        setHasCustomSplit(false);
        onSubmit();
      }}
      className=""
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
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="payer">
          Payer
        </label>
        <select
          {...register('payer', {
            required: 'Please select a payer',
          })}
          className="mt-2 block rounded bg-white bg-clip-padding bg-no-repeat px-3 py-2 text-base font-normal text-gray-700 outline outline-1 transition ease-in-out focus:border-indigo-600 focus:bg-white focus:text-gray-700"
        >
          {tracker.users.map((userId) => (
            <option key={userId} value={userId}>
              {userId}
            </option>
          ))}
        </select>
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
            <optgroup key={mainKey} label={mainCategories[mainKey].en_US}>
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
      <div className="mt-5 inline-flex items-center">
        <input
          {...register('settled')}
          type="checkbox"
          id="settled"
          className="mr-2"
        />
        <label
          className="font-semibold lowercase text-slate-600"
          htmlFor="settled"
        >
          Transaction has been settled
        </label>
      </div>
      <div className="mt-5">
        <input
          type="checkbox"
          checked={hasCustomSplit}
          onChange={(e) => {
            setHasCustomSplit(e.target.checked);
          }}
          className="mr-2"
          id="customSplit"
        />
        <label
          className="font-semibold lowercase text-slate-600"
          htmlFor="customSplit"
        >
          Custom split
        </label>
      </div>
      {hasCustomSplit && (
        <div className="mt-2">
          <input
            {...register('split')}
            className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 placeholder-slate-400 focus:ring-0"
            type="text"
          />
        </div>
      )}
      {children}
    </form>
  );
}
