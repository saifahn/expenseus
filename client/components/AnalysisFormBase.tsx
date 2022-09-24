import { ChangeEvent } from 'react';
import { UseFormRegister, UseFormSetValue } from 'react-hook-form';
import { dateRanges } from 'utils/dates';

export type AnalysisFormInputs = {
  from: string;
  to: string;
};

type Props = {
  register: UseFormRegister<AnalysisFormInputs>;
  onSubmit: () => void;
  setValue: UseFormSetValue<AnalysisFormInputs>;
};

export default function AnalysisFormBase({
  register,
  onSubmit,
  setValue,
}: Props) {
  function handlePresetSelect(e: ChangeEvent<HTMLSelectElement>) {
    const preset = e.target.value;
    const { from, to } = dateRanges[preset].presetFn();
    setValue('from', from);
    setValue('to', to);
  }

  return (
    <form
      className="bg-white py-3"
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit();
      }}
    >
      <h3 className="text-lg font-bold lowercase">Analyze transactions</h3>
      <div className="mt-3">
        <label className="block font-semibold lowercase text-slate-600">
          date preset
        </label>
        <select
          className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 text-center lowercase placeholder-slate-400 focus:ring-0"
          onChange={handlePresetSelect}
        >
          {Object.entries(dateRanges).map(([preset, { name }]) => (
            <option key={preset} value={preset}>
              {name}
            </option>
          ))}
        </select>
      </div>
      <div className="mt-5">
        <label
          className="block font-semibold lowercase text-slate-600"
          htmlFor="dateFrom"
        >
          From
        </label>
        <input
          {...register('from', { required: 'Please input a date' })}
          className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 text-center placeholder-slate-400 focus:ring-0"
          type="date"
          id="dateFrom"
        />
      </div>
      <div className="mt-5">
        <label
          className="block font-semibold lowercase text-slate-600"
          htmlFor="dateTo"
        >
          To
        </label>
        <input
          {...register('to', { required: 'Please input a date' })}
          className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 text-center placeholder-slate-400 focus:ring-0"
          type="date"
          id="dateTo"
        />
      </div>

      <div className="mt-4 flex justify-end">
        <button
          className="rounded bg-violet-500 py-2 px-4 font-medium lowercase text-white hover:bg-violet-700 focus:outline-none focus:ring"
          type="submit"
        >
          Analyze
        </button>
      </div>
    </form>
  );
}
