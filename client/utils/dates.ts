import { Temporal } from 'temporal-polyfill';

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

// this month
// last month
// this week - needs start of week
// last week
// last 7 days
// last month
// last 90 days
// last 180 days
export const presets = {
  now() {
    return Temporal.Now.plainDateISO();
  },
  startOfWeek() {
    return presets.now().subtract({ days: presets.now().dayOfWeek - 1 });
  },
};
