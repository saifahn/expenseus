import { UserAlreadyExistsError } from 'ddb/errors';
import { setUpUserRepo } from 'ddb/setUpRepos';
import { AuthOptions } from 'next-auth';
import NextAuth from 'next-auth/next';
import GoogleProvider from 'next-auth/providers/google';

export const authOptions: AuthOptions = {
  secret: 'test-secret',
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
      wellKnown: 'https://accounts.google.com/.well-known/openid-configuration',
    }),
  ],
  callbacks: {
    async signIn({ user }) {
      const userRepo = setUpUserRepo();
      try {
        await userRepo.createUser({
          id: user.email!,
          username: user.email!,
          name: user.name!,
        });
        console.log('user successfully created');
        return true;
      } catch (err) {
        if (err instanceof UserAlreadyExistsError) {
          return true;
        }
        console.error(err);
        return false;
      }
    },
  },
};

export default NextAuth(authOptions);
