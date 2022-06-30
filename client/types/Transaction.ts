import { CategoryKey } from 'data/categories';

export interface Transaction {
  id: string;
  userId: string;
  location: string;
  amount: number;
  imageUrl?: string;
  date: number;
  category: CategoryKey;
  details: string;
}
