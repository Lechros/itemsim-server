import { createCors } from "itty-router";

export const { preflight, corsify } = createCors({
  methods: ["GET"],
  origins: ["https://itemsim.pages.dev", "http://localhost:5173"],
});
