const categories = {
  'unspecified.unspecified': {
    en_US: 'Unspecified',
  },
  'food.food': {
    en_US: 'Food',
  },
  'food.groceries': {
    en_US: 'Groceries',
  },
  'food.eating-out': {
    en_US: 'Eating out',
  },
  'food.delivery': {
    en_US: 'Delivery',
  },
  'daily.daily': {
    en_US: 'Daily goods',
  },
  'daily.consumables': {
    en_US: 'Consumables',
  },
  'daily.children': {
    en_US: 'Child-related',
  },
  'transport.transport': {
    en_US: 'Transportation',
  },
  'transport.train': {
    en_US: 'Train',
  },
  'transport.bus': {
    en_US: 'Bus',
  },
  'transport.taxi': {
    en_US: 'Taxi',
  },
  'utilities.utilities': {
    en_US: 'Utilities',
  },
  'utilities.electricity': {
    en_US: 'Electricity',
  },
  'utilities.water': {
    en_US: 'Water',
  },
  'utilities.gas': {
    en_US: 'Gas',
  },
  'utilities.internet': {
    en_US: 'Internet',
  },
  'utilities.phone': {
    en_US: 'Phone',
  },
  'entertainment.entertainment': {
    en_US: 'Entertainment',
  },
  'entertainment.film-video': {
    en_US: 'Film and Video',
  },
  'entertainment.books': {
    en_US: 'Books',
  },
  'entertainment.music': {
    en_US: 'Music',
  },
  'entertainment.leisure': {
    en_US: 'Leisure',
  },
  'entertainment.hobbies': {
    en_US: 'Hobbies',
  },
  'entertainment.subscriptions': {
    en_US: 'Subscriptions',
  },
  'clothing.clothing': {
    en_US: 'Clothing',
  },
  'clothing.footwear': {
    en_US: 'Footwear',
  },
  'beauty.beauty': {
    en_US: 'Beauty',
  },
  'beauty.hair': {
    en_US: 'Hair',
  },
  'beauty.cosmetics': {
    en_US: 'Cosmetics',
  },
  'travel.travel': {
    en_US: 'Travel',
  },
  'travel.transport': {
    en_US: 'Transport',
  },
  'travel.accommodation': {
    en_US: 'Accommodation',
  },
  'home.home': {
    en_US: 'Home',
  },
  'home.furniture': {
    en_US: 'Furniture',
  },
  'home.electronics': {
    en_US: 'Electronics',
  },
  'home.kitchen': {
    en_US: 'Kitchen',
  },
  'medical.medical': {
    en_US: 'Medical',
  },
  'other.other': {
    en_US: 'Other',
  },
  'social.social': {
    en_US: 'Social',
  },
  'education.education': {
    en_US: 'Education',
  },
  'education.lessons': {
    en_US: 'Lessons',
  },
  'housing.housing': {
    en_US: 'Housing',
  },
  'housing.rent': {
    en_US: 'Rent',
  },
  'housing.mortgage': {
    en_US: 'Mortgage',
  },
  'insurance.insurance': {
    en_US: 'Insurance',
  },
  'car.car': {
    en_US: 'Car',
  },
};

export type CategoryKey = keyof typeof categories;

const categoryKeys = Object.keys(categories) as CategoryKey[];

export const enUSCategories = categoryKeys.map((key) => ({
  key,
  value: categories[key].en_US,
}));

export default categories;
