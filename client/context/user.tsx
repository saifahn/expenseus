import { User } from 'components/UserList';
import { HttpException } from 'config/fetcher';
import { createContext, PropsWithChildren, useContext } from 'react';
import useSWR from 'swr';

interface UserContextState {
  user?: User;
  error?: HttpException;
}

const defaultState: UserContextState = {
  user: undefined,
  error: undefined,
};

const UserContext = createContext<UserContextState>(defaultState);

export function UserProvider({ children }: PropsWithChildren<{}>) {
  const { data: user, error } = useSWR<User, HttpException>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/users/self`,
  );

  return (
    <UserContext.Provider value={{ user, error }}>
      {children}
    </UserContext.Provider>
  );
}

export function useUserContext() {
  return useContext(UserContext);
}
