/**
 * Welcome to Cloudflare Workers! This is your first worker.
 *
 * - Run `wrangler dev src/index.ts` in your terminal to start a development server
 * - Open a browser tab at http://localhost:8787/ to see your worker in action
 * - Run `wrangler publish src/index.ts --name my-worker` to publish your worker
 *
 * Learn more at https://developers.cloudflare.com/workers/
 */

export interface Env {
  // Example binding to KV. Learn more at https://developers.cloudflare.com/workers/runtime-apis/kv/
  // MY_KV_NAMESPACE: KVNamespace;
  //
  // Example binding to Durable Object. Learn more at https://developers.cloudflare.com/workers/runtime-apis/durable-objects/
  // MY_DURABLE_OBJECT: DurableObjectNamespace;
  //
  // Example binding to R2. Learn more at https://developers.cloudflare.com/workers/runtime-apis/r2/
  // MY_BUCKET: R2Bucket;
  //
  // Example binding to a Service. Learn more at https://developers.cloudflare.com/workers/runtime-apis/service-bindings/
  // MY_SERVICE: Fetcher;
}

import {
  GearData,
  createGearFromNode,
  createPotentialFromCode,
  gearJson,
} from "@malib/create-gear";
import { gearToPlain } from "@malib/gear";

const gearDatas = Object.entries(gearJson);

function search(names: [string, GearData][], word: string) {
  if (word.length < 1) {
    return [];
  }
  return names
    .filter((e) => match(e[1].name, word))
    .sort((a, b) => compare(a[1].name, b[1].name, word));
}

function match(haystack: string, word: string) {
  let j = 0;
  for (let i = 0; i < haystack.length && j < word.length; i++) {
    if (haystack[i] === word[j]) {
      j++;
    }
  }
  return j === word.length;
}

function compare(name1: string, name2: string, word: string) {
  // name is same as word, preserve itemID ordering
  if (name1 == word && name2 == word) return 0;
  if (name1 == word) return -1;
  if (name2 == word) return 1;
  const index1 = name1.indexOf(word);
  const index2 = name2.indexOf(word);
  // both names contain word
  if (index1 >= 0 && index2 >= 0) {
    // both have same word position
    if (index1 == index2) return 0;
    // word appearing earlier comes front
    return index1 - index2;
  }
  // only one name contains word
  if (index1 >= 0) return -1;
  if (index2 >= 0) return 1;
  // both are only fuzzy-matched, preserve itemID ordering
  return 0;
}

export default {
  async fetch(
    request: Request,
    env: Env,
    ctx: ExecutionContext
  ): Promise<Response> {
    const input = /https:\/\/.+?\.workers\.dev\/(.+?)$/.exec(request.url)?.[1];
    if (!input)
      return Response.json(
        {},
        {
          status: 400,
          headers: {
            "Access-Control-Allow-Credentials": "true",
            "Access-Control-Allow-Methods": "GET",
            "Access-Control-Allow-Origin": "https://itemsim.pages.dev",
            "Access-Control-Allow-Headers": "Content-Type",
          },
        }
      );
    const params = new URLSearchParams(input);
    if (params.has("query")) {
      const query = params.get("query") as string;

      const items = search(gearDatas, query).map((e) => [
        Number(e[0]),
        gearToPlain(
          createGearFromNode(e[1], Number(e[0]), createPotentialFromCode)
        ),
      ]);

      return Response.json(items, {
        status: 200,
        headers: {
          "Access-Control-Allow-Credentials": "true",
          "Access-Control-Allow-Methods": "GET",
          "Access-Control-Allow-Origin": "https://itemsim.pages.dev",
          "Access-Control-Allow-Headers": "Content-Type",
        },
      });
    }

    return Response.json(
      {},
      {
        status: 400,
        headers: {
          "Access-Control-Allow-Credentials": "true",
          "Access-Control-Allow-Methods": "GET",
          "Access-Control-Allow-Origin": "https://itemsim.pages.dev",
          "Access-Control-Allow-Headers": "Content-Type",
        },
      }
    );
  },
};
