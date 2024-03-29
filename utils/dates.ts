import { Temporal } from 'temporal-polyfill';

const monthNameMapEN = {
  1: 'January',
  2: 'February',
  3: 'March',
  4: 'April',
  5: 'May',
  6: 'June',
  7: 'July',
  8: 'August',
  9: 'September',
  10: 'October',
  11: 'November',
  12: 'December',
} as const;

export type MonthEN = (typeof monthNameMapEN)[keyof typeof monthNameMapEN];

/**
 * Takes a date in epoch seconds format and returns its year in the UTC
 * timezone.
 */
export function epochSecToUTCYear(date: number) {
  return Temporal.Instant.fromEpochSeconds(date).toZonedDateTimeISO('UTC').year;
}

/**
 * Takes a date in epoch seconds format and converts it to a month name
 * in English based on UTC timezone.
 */
export function epochSecToUTCMonthEN(date: number): MonthEN {
  return monthNameMapEN[
    Temporal.Instant.fromEpochSeconds(date).toZonedDateTimeISO('UTC')
      .month as keyof typeof monthNameMapEN
  ];
}

/**
 * Takes a date in epoch seconds format and converts it to a locale string
 * to be displayed.
 */
export function epochSecToLocaleString(date: number) {
  return Temporal.Instant.fromEpochSeconds(date)
    .toZonedDateTimeISO('UTC')
    .toPlainDate()
    .toLocaleString();
}

/**
 * Takes a date in epoch seconds format and converts it to a ISO date string to
 * be used in the date input.
 */
export function epochSecToISOString(date: number) {
  return Temporal.Instant.fromEpochSeconds(date)
    .toZonedDateTimeISO('UTC')
    .toPlainDate()
    .toString();
}

/**
 * Gets the current date now in ISO format for use in e.g. a date input.
 */
export function plainDateISONowString() {
  return Temporal.Now.plainDateISO().toString();
}

/**
 * Takes a plainDate ISO format string and converts it to epoch seconds to be
 * submitted to the back end.
 */
export function plainDateStringToEpochSec(date: string) {
  // z to use the UTC timezone for all dates submitted
  return Temporal.Instant.from(`${date.toString()}z`).epochSeconds;
}

export const presets = {
  now() {
    return Temporal.Now.plainDateISO();
  },
  startOfWeek() {
    return presets.now().subtract({ days: presets.now().dayOfWeek - 1 });
  },
  startOfLastWeek() {
    return presets
      .now()
      .subtract({ weeks: 1, days: presets.now().dayOfWeek - 1 });
  },
  endOfLastWeek() {
    return presets.now().subtract({ days: presets.now().dayOfWeek });
  },
  startOfMonth() {
    return presets.now().subtract({ days: presets.now().day - 1 });
  },
  startOfLastMonth() {
    return presets.now().subtract({ months: 1, days: presets.now().day - 1 });
  },
  endOfLastMonth() {
    return presets.now().subtract({ days: presets.now().day });
  },
  sevenDaysAgo() {
    return presets.now().subtract({ days: 7 });
  },
  thirtyDaysAgo() {
    return presets.now().subtract({ days: 30 });
  },
  ninetyDaysAgo() {
    return presets.now().subtract({ days: 90 });
  },
  oneHundredAndEightyDaysAgo() {
    return presets.now().subtract({ days: 180 });
  },
};

export type DateRangePresetFn = () => { from: string; to: string };
export type DateRangePresets = {
  [key: string]: {
    name: string;
    presetFn: DateRangePresetFn;
  };
};

export const dateRanges: DateRangePresets = {
  thisWeek: {
    name: 'This week',
    presetFn: function () {
      return {
        from: presets.startOfWeek().toString(),
        to: presets.now().toString(),
      };
    },
  },
  lastWeek: {
    name: 'Last week',
    presetFn: function () {
      return {
        from: presets.startOfLastWeek().toString(),
        to: presets.endOfLastWeek().toString(),
      };
    },
  },
  thisMonth: {
    name: 'This month',
    presetFn: function () {
      return {
        from: presets.startOfMonth().toString(),
        to: presets.now().toString(),
      };
    },
  },
  lastMonth: {
    name: 'Last month',
    presetFn: function () {
      return {
        from: presets.startOfLastMonth().toString(),
        to: presets.endOfLastMonth().toString(),
      };
    },
  },
  lastNinetyDays: {
    name: 'Last 90 days',
    presetFn: function () {
      return {
        from: presets.ninetyDaysAgo().toString(),
        to: presets.now().toString(),
      };
    },
  },
};

export function formatDateForTxnCard(date: number) {
  const currentYear = Temporal.Now.zonedDateTimeISO('UTC').year;
  // the date is stored as a epoch seconds, Date constructor takes milliseconds
  return new Date(date * 1000).toLocaleDateString(['en-GB', 'ja-JP'], {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    // should be able to assume that dates without a year are from current year
    ...(epochSecToUTCYear(date) !== currentYear && {
      year: 'numeric',
    }),
  });
}
