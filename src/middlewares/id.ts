import { IRequest } from "itty-router";

export const withId = (request: IRequest) => {
  const { id: idStr } = request.params;
  const id = Number.parseInt(idStr);
  if (!isNaN(id)) {
    request.id = id;
  }
};
