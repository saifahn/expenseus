import "@testing-library/jest-dom";
import fetch from "cross-fetch";
import dotenv from "dotenv";

// use ponyfill in tests because fetch is not in node
global.fetch = fetch;

dotenv.config({ path: ".env.test" });
