import typing

from flask import Flask

from generated.server import Server
import generated.message as message

app = Flask(__name__)

server = Server(app)


@server.testTag1.op1
async def op1(
    uri: message.testTag1.op1.uri,
    qry: message.testTag1.op1.qry,
    req: message.testTag1.op1.req,
) -> typing.Tuple[typing.Optional[message.testTag1.op1.rsp200]]:

    rsp200 = message.testTag1.op1.rsp200(
        uri1=uri.uri1,
        uri2=uri.uri2,
        qry1=qry.qry1,
        qry2=qry.qry2,
        qryo=qry.qryo,
        req1=req.req1,
        req2=req.req2,
    )

    return (rsp200,)


@server.testTag1.op2
async def op2(
    uri: message.testTag1.op2.uri,
    qry: message.testTag1.op2.qry,
    req: message.testTag1.op2.req,
) -> typing.Tuple[typing.Optional[message.testTag1.op2.rsp200]]:
    rsp200 = req

    return (rsp200,)


@server.testTag1.op3
async def op3(
    uri: message.testTag1.op3.uri,
    qry: message.testTag1.op3.qry,
    req: message.testTag1.op3.req,
) -> typing.Tuple[typing.Optional[message.testTag1.op3.rsp200]]:
    rsp200 = req

    return (rsp200,)


@server.testTag2.op4
async def op4() -> typing.Tuple[typing.Optional[message.testTag2.op4.rsp200]]:
    rsp200 = ""
    return (rsp200,)


@server.testTag2.op5
async def op5() -> typing.Tuple[typing.Optional[message.testTag2.op5.rsp204]]:
    rsp200 = ""
    return (rsp200,)


@server.testTag2.op6
async def op6(
    req: message.testTag2.op6.req,
) -> typing.Tuple[
    typing.Optional[message.testTag2.op6.rsp200],
    typing.Optional[message.testTag2.op6.rsp201],
    typing.Optional[message.testTag2.op6.rsp202],
    typing.Optional[message.testTag2.op6.rsp203],
]:
    rsp200 = None
    rsp201 = None
    rsp202 = None
    rsp203 = None

    if req.code == 200:
        rsp200 = message.testTag2.op6.rsp200()
    if req.code == 201:
        rsp201 = message.testTag2.op6.rsp201()
    if req.code == 202:
        rsp202 = message.testTag2.op6.rsp202()
    if req.code == 203:
        rsp203 = message.testTag2.op6.rsp203()
    return rsp200, rsp201, rsp202, rsp203


if __name__ == "__main__":
    app.run(host="127.0.0.1", port=30435)
