import { UserAPI } from 'api';
import { User } from 'components/UserList';
import { createContext, useContext, useEffect, useRef, useState } from 'react';

interface UserFetchStatus {
  status: 'idle' | 'fulfilled' | 'rejected' | 'loading';
  error: number;
}

interface UserContextState {
  user: User;
  userFetchStatus: UserFetchStatus;
}

const defaultState: UserContextState = {
  user: {
    username: null,
    name: null,
    id: null,
  },
  userFetchStatus: {
    status: 'idle',
    error: null,
  },
};

const UserContext = createContext<UserContextState>(defaultState);

export function UserProvider({ children }) {
  const [user, setUser] = useState(defaultState.user);
  const [{ status, error }, setStatus] = useState<UserFetchStatus>(
    defaultState.userFetchStatus,
  );
  const state: UserContextState = {
    user,
    userFetchStatus: {
      status,
      error,
    },
  };

  const cancelled = useRef(false);

  async function fetchSelf() {
    try {
      setStatus({ status: 'loading', error: null });
      const api = new UserAPI();
      const loggedInUser = await api.getSelf();
      if (!cancelled.current) {
        setUser(loggedInUser);
        setStatus({ status: 'fulfilled', error: null });
      }
    } catch (error) {
      if (!cancelled.current) {
        return setStatus({ status: 'rejected', error });
      }
    }
  }

  useEffect(() => {
    fetchSelf();
    return () => {
      cancelled.current = true;
    };
  }, []);

  return <UserContext.Provider value={state}>{children}</UserContext.Provider>;
}

export function useUserContext() {
  return useContext(UserContext);
}
