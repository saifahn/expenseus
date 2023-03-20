import { z } from 'zod';

export const mainCategories = {
  unspecified: {
    en_US: 'Unspecified',
    emoji: '🏷',
    colour: '#a7f3d0', // emerald 200
  },
  food: {
    en_US: 'Food & Drink',
    emoji: '🍻',
    colour: '#ddd6fe', // violet 200
  },
  daily: {
    en_US: 'Daily',
    emoji: '🧽',
    colour: '#fecdd3', // rose 200
  },
  transport: {
    en_US: 'Transport',
    emoji: '🚆',
    colour: '#6ee7b7', // emerald 300
  },
  utilities: {
    en_US: 'Utilities',
    emoji: '💡',
    colour: '#c4b5fd', // violet 300
  },
  entertainment: {
    en_US: 'Entertainment',
    emoji: '🎥',
    colour: '#fda4af', // rose 300
  },
  clothing: {
    en_US: 'Clothing',
    emoji: '👚',
    colour: '#34d399', // emerald 400
  },
  beauty: {
    en_US: 'Beauty',
    emoji: '💇',
    colour: '#a78bfa', // violet 400
  },
  travel: {
    en_US: 'Travel',
    emoji: '✈️',
    colour: '#fb7185', // rose 400
  },
  home: {
    en_US: 'Home',
    emoji: '🍴',
    colour: '#10b981', // emerald 500
  },
  medical: {
    en_US: 'Medical',
    emoji: '🩺',
    colour: '#8b5cf6', // violet 500
  },
  social: {
    en_US: 'Social',
    emoji: '🫂',
    colour: '#f43f5e', // rose 500
  },
  education: {
    en_US: 'Education',
    emoji: '🎓',
    colour: '#059669', // emerald 600
  },
  housing: {
    en_US: 'Housing',
    emoji: '🏠',
    colour: '#7c3aed', // violet 600
  },
  insurance: {
    en_US: 'Insurance',
    emoji: '🛡️',
    colour: '#e11d48', // rose 600
  },
  car: {
    en_US: 'Car',
    emoji: '🚗',
    colour: '#047857', // emerald 700
  },
  other: {
    en_US: 'Other',
    emoji: '🔸',
    colour: '#6d28d9', // violet 700
  },
};

export function getEmojiForTxnCard(key: SubcategoryKey) {
  const mc = subcategories[key].mainCategory;
  return mainCategories[mc].emoji;
}

export const categoryColours = Object.values(mainCategories).map(
  (item) => item.colour,
);

export type MainCategoryKey = keyof typeof mainCategories;

export const mainCategoryKeys = Object.keys(
  mainCategories,
) as MainCategoryKey[];

const subcategoryKeys = [
  'unspecified.unspecified',
  'food.food',
  'food.groceries',
  'food.eating-out',
  'food.delivery',
  'daily.daily',
  'daily.consumables',
  'daily.children',
  'transport.transport',
  'transport.train',
  'transport.bus',
  'transport.taxi',
  'utilities.utilities',
  'utilities.electricity',
  'utilities.gas',
  'utilities.water',
  'utilities.internet',
  'utilities.phone',
  'entertainment.entertainment',
  'entertainment.books',
  'entertainment.leisure',
  'entertainment.hobbies',
  'entertainment.music',
  'entertainment.film-video',
  'entertainment.subscriptions',
  'clothing.clothing',
  'clothing.footwear',
  'beauty.beauty',
  'beauty.cosmetics',
  'beauty.hair',
  'travel.travel',
  'travel.transport',
  'travel.accommodation',
  'home.home',
  'home.electronics',
  'home.furniture',
  'home.kitchen',
  'medical.medical',
  'other.other',
  'social.social',
  'social.gifts',
  'education.education',
  'education.lessons',
  'housing.housing',
  'housing.rent',
  'housing.mortgage',
  'car.car',
  'insurance.insurance',
] as const;
export const SubcategoryKeys = z.enum(subcategoryKeys);
export type SubcategoryKey = z.infer<typeof SubcategoryKeys>;

// use when satisfies is supported?
type SubCategories = {
  [k in SubcategoryKey]: {
    mainCategory: MainCategoryKey;
    en_US: string;
  };
};

export const subcategories: SubCategories = {
  'unspecified.unspecified': {
    mainCategory: 'unspecified',
    en_US: 'Unspecified',
  },
  'food.food': {
    mainCategory: 'food',
    en_US: 'Food & Drink',
  },
  'food.groceries': {
    mainCategory: 'food',
    en_US: 'Groceries',
  },
  'food.eating-out': {
    mainCategory: 'food',
    en_US: 'Eating out',
  },
  'food.delivery': {
    mainCategory: 'food',
    en_US: 'Delivery',
  },
  'daily.daily': {
    mainCategory: 'daily',
    en_US: 'Daily goods',
  },
  'daily.consumables': {
    mainCategory: 'daily',
    en_US: 'Consumables',
  },
  'daily.children': {
    mainCategory: 'daily',
    en_US: 'Child-related',
  },
  'transport.transport': {
    mainCategory: 'transport',
    en_US: 'Transportation',
  },
  'transport.train': {
    mainCategory: 'transport',
    en_US: 'Train',
  },
  'transport.bus': {
    mainCategory: 'transport',
    en_US: 'Bus',
  },
  'transport.taxi': {
    mainCategory: 'transport',
    en_US: 'Taxi',
  },
  'utilities.utilities': {
    mainCategory: 'utilities',
    en_US: 'Utilities',
  },
  'utilities.electricity': {
    mainCategory: 'utilities',
    en_US: 'Electricity',
  },
  'utilities.water': {
    mainCategory: 'utilities',
    en_US: 'Water',
  },
  'utilities.gas': {
    mainCategory: 'utilities',
    en_US: 'Gas',
  },
  'utilities.internet': {
    mainCategory: 'utilities',
    en_US: 'Internet',
  },
  'utilities.phone': {
    mainCategory: 'utilities',
    en_US: 'Phone',
  },
  'entertainment.entertainment': {
    mainCategory: 'entertainment',
    en_US: 'Entertainment',
  },
  'entertainment.film-video': {
    mainCategory: 'entertainment',
    en_US: 'Film and Video',
  },
  'entertainment.books': {
    mainCategory: 'entertainment',
    en_US: 'Books',
  },
  'entertainment.music': {
    mainCategory: 'entertainment',
    en_US: 'Music',
  },
  'entertainment.leisure': {
    mainCategory: 'entertainment',
    en_US: 'Leisure',
  },
  'entertainment.hobbies': {
    mainCategory: 'entertainment',
    en_US: 'Hobbies',
  },
  'entertainment.subscriptions': {
    mainCategory: 'entertainment',
    en_US: 'Subscriptions',
  },
  'clothing.clothing': {
    mainCategory: 'clothing',
    en_US: 'Clothing',
  },
  'clothing.footwear': {
    mainCategory: 'clothing',
    en_US: 'Footwear',
  },
  'beauty.beauty': {
    mainCategory: 'beauty',
    en_US: 'Beauty',
  },
  'beauty.hair': {
    mainCategory: 'beauty',
    en_US: 'Hair',
  },
  'beauty.cosmetics': {
    mainCategory: 'beauty',
    en_US: 'Cosmetics',
  },
  'travel.travel': {
    mainCategory: 'travel',
    en_US: 'Travel',
  },
  'travel.transport': {
    mainCategory: 'travel',
    en_US: 'Transport',
  },
  'travel.accommodation': {
    mainCategory: 'travel',
    en_US: 'Accommodation',
  },
  'home.home': {
    mainCategory: 'home',
    en_US: 'Home',
  },
  'home.furniture': {
    mainCategory: 'home',
    en_US: 'Furniture',
  },
  'home.electronics': {
    mainCategory: 'home',
    en_US: 'Electronics',
  },
  'home.kitchen': {
    mainCategory: 'home',
    en_US: 'Kitchen',
  },
  'medical.medical': {
    mainCategory: 'medical',
    en_US: 'Medical',
  },
  'social.social': {
    mainCategory: 'social',
    en_US: 'Social',
  },
  'social.gifts': {
    mainCategory: 'social',
    en_US: 'Gifts',
  },
  'education.education': {
    mainCategory: 'education',
    en_US: 'Education',
  },
  'education.lessons': {
    mainCategory: 'education',
    en_US: 'Lessons',
  },
  'housing.housing': {
    mainCategory: 'housing',
    en_US: 'Housing',
  },
  'housing.rent': {
    mainCategory: 'housing',
    en_US: 'Rent',
  },
  'housing.mortgage': {
    mainCategory: 'housing',
    en_US: 'Mortgage',
  },
  'insurance.insurance': {
    mainCategory: 'insurance',
    en_US: 'Insurance',
  },
  'car.car': {
    mainCategory: 'car',
    en_US: 'Car',
  },
  'other.other': {
    mainCategory: 'other',
    en_US: 'Other',
  },
};

export const categoryKeys = Object.keys(subcategories) as SubcategoryKey[];

type AllCategories = {
  [mainCat in MainCategoryKey]: Array<SubcategoryKey>;
};

// FIXME: this is kind of awful and hard to understand
/**
 * makes an AllCategories from the subcategories and main categories for making
 * tiered options with `optgroup`
 */
export const categories = categoryKeys.reduce((prev, key) => {
  const mainCategory = subcategories[key].mainCategory;
  if (!prev[mainCategory]) {
    prev[mainCategory] = [];
  }
  prev[mainCategory].push(key);
  return prev;
}, {} as AllCategories);

export const categoryNameFromKeyEN = (key: SubcategoryKey) =>
  subcategories[key].en_US;
