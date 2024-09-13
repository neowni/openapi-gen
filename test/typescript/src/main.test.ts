import { test, expect } from "vitest";

import axios from "axios";
import client from "./generated/client";
import * as message from "./generated/message";
import * as models from "./generated/models";

test("test", async () => {
  const c = new client(
    axios.create({
      baseURL: "http://127.0.0.1:30435/",
    })
  );

  function randomEnum() {
    switch (randomInt() % 3) {
      case 0:
        return models.Enum.value1;
      case 1:
        return models.Enum.value2;
      default:
        return models.Enum.value3;
    }
  }

  //                                                                            op1
  {
    const rspE: message.testTag1.op1Rsp200 = {
      uri1: randomString(),
      uri2: randomInt(),
      qry1: randomString(),
      qry2: randomInt(),
      req1: randomInt(),
      req2: Array.from({ length: 16 }, randomString),
    };

    const rspA = await c.testTag1.op1(
      {
        uri1: rspE.uri1,
        uri2: rspE.uri2,
      },
      {
        qry1: rspE.qry1,
        qry2: rspE.qry2,
      },
      {
        req1: rspE.req1,
        req2: rspE.req2,
      }
    );

    expect(rspA._200).toStrictEqual(rspE);
  }

  //                                                                            op2
  {
    const rdObj2 = () => {
      const o: models.Object2 = {
        requiredField: randomString(),
      };

      if (randomBoolean()) {
        o.optionalField = randomString();
      }

      return o;
    };

    const rspE: message.testTag1.op2Rsp200 = {
      stringField: randomString(),
      intField: randomInt(),
      floatField: randomFloat(),
      enumField: randomEnum(),
      arrayField1: Array.from({ length: 16 }, randomInt),
      arrayField2: Array.from({ length: 16 }, rdObj2),
      objectField1: {},
      objectField2: rdObj2(),
    };

    const rspA = await c.testTag1.op2(
      {
        uri1: randomString(),
        uri2: randomInt(),
      },
      {
        qry1: randomString(),
      },
      rspE
    );

    expect(rspA._200).toStrictEqual(rspE);
  }

  //                                                                            op3
  {
    const rspE = randomString();

    const rspA = await c.testTag1.op3(
      {
        uri1: randomString(),
        uri2: randomInt(),
      },
      {
        qry1: randomString(),
      },
      rspE
    );

    expect(rspA._200).toStrictEqual(rspE);
  }

  //                                                                            op4
  {
    const rspA = await c.testTag2.op4();

    expect(rspA._200).not.toBeUndefined();
  }

  //                                                                            op5
  {
    const rspA = await c.testTag2.op5();

    expect(rspA._204).not.toBeUndefined();
  }

  //                                                                            op6
  {
    const code = Math.floor(Math.random() * 4) + 200;

    const rspA = await c.testTag2.op6({
      code: code,
    });

    switch (code) {
      case 200:
        expect(rspA._200).not.toBeUndefined();
        break;

      case 201:
        expect(rspA._201).not.toBeUndefined();
        break;

      case 202:
        expect(rspA._202).not.toBeUndefined();
        break;

      case 203:
        expect(rspA._203).not.toBeUndefined();
        break;
    }
  }
});

function randomString(): string {
  const characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  let result = "";
  const charactersLength = characters.length;

  for (let i = 0; i < 16; i++) {
    const randomIndex = Math.floor(Math.random() * charactersLength);
    result += characters.charAt(randomIndex);
  }

  return result;
}

function randomInt(): number {
  return Math.floor(Math.random() * 4294967296);
}

function randomFloat(): number {
  return Math.random();
}

function randomBoolean(): boolean {
  return Math.random() > 0.5;
}
