import {
  categories,
  SubcategoryKey,
  mainCategories,
  mainCategoryKeys,
  subcategories,
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
        onSubmit();
      }}
      className="border-4 p-6"
    >
      {title && <h3 className="text-lg font-semibold">{title}</h3>}
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="location">
          Location
        </label>
        <input
          {...register('location', {
            required: 'Please input a location',
          })}
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
          type="text"
          id="location"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="amount">
          Amount
        </label>
        <input
          {...register('amount', {
            min: {
              value: 1,
              message: 'Please input a positive amount',
            },
            required: 'Please input an amount',
          })}
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
          {...register('date', {
            required: 'Please input a date',
          })}
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
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="settled">
          Settled?
        </label>
        <input {...register('settled')} type="checkbox" id="settled" />
      </div>
      <div className="mt-4">
        <label className="block font-semibold">Category</label>
        <select
          {...register('category')}
          className="mt-2 block rounded bg-white bg-clip-padding bg-no-repeat px-3 py-2 text-base font-normal text-gray-700 outline outline-1 transition ease-in-out focus:border-indigo-600 focus:bg-white focus:text-gray-700"
        >
          {mainCategoryKeys.map((mainKey) => (
            <optgroup key={mainKey} label={mainCategories[mainKey].en_US}>
              {categories[mainKey].map((subKey) => (
                <option key={subKey} value={subKey}>
                  {subcategories[subKey].en_US}
                </option>
              ))}
            </optgroup>
          ))}
        </select>
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="details">
          Details
        </label>
        <input
          {...register('details')}
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
          type="text"
          id="details"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="details">
          Custom split?
        </label>
        <input
          type="checkbox"
          checked={hasCustomSplit}
          onChange={(e) => {
            setHasCustomSplit(e.target.checked);
          }}
        />
      </div>
      {hasCustomSplit && (
        <div className="mt-4">
          <input
            {...register('split')}
            className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
            type="text"
          />
        </div>
      )}
      {children}
    </form>
  );
}
