import {
  categories,
  SubcategoryKey,
  mainCategories,
  mainCategoryKeys,
  categoryNameFromKeyEN,
} from 'data/categories';
import { CreateSharedTxnPayload } from 'pages/api/v1/trackers/[trackerId]/transactions';
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

export function makeSharedTxnPayload(
  data: SharedTxnFormInputs,
): Omit<CreateSharedTxnPayload, 'tracker'> {
  const participants = data.participants.split(',');

  let split;
  if (data.split) {
    // the format is `userid:split,userid:split` originally, which can be split
    const [first, second] = data.split.split(',');
    const [firstUser, firstSplit] = first.split(':');
    const [secondUser, secondSplit] = second.split(':');
    split = {
      [firstUser]: Number(firstSplit),
      [secondUser]: Number(secondSplit),
    };
  }

  const payload = {
    ...data,
    date: plainDateStringToEpochSec(data.date),
    amount: Number(data.amount),
    participants,
    split,
    ...(data.settled && { settled: data.settled }),
  };

  return payload;
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
      className="bg-white py-3"
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
        <label
          className="block font-semibold lowercase text-slate-600"
          htmlFor="payer"
        >
          Payer
        </label>
        <select
          {...register('payer', {
            required: 'Please select a payer',
          })}
          className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 lowercase placeholder-slate-400 focus:ring-0"
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
