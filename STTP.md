# Simple Text Transfer Protocol (STTP)

## Content
- [Content](#content)
- [What is it?](#what-is-it)
- [Overview](#overview)
- [Structure](#structure)
- [URL](#url)

## What is it?

It is an application protocol created for learning purposes only. Don't intended to be real-world protocol. So, take it easy.

## Overview

The Simple Text Transfer Protocol (STTP) is an application layer protocol for transmitting text data.

STTP can be used both as request-response protocol or as one-way protocol. It is not intended nor for P2P architecture nor for client-server model. It is intended for direct connections between two clients, primarily in local area network.

STTP requires an underlying and reliable transport layer protocol. TCP, for example.

STTP allows processing of transmitted data by different handlers at application level. Every handler may have different logic to handle transmitted data. For routing concept of ports is used.

## Structure

It is a structure of bytes that are comes from transport protocol.

<table>
  <tr align="center">
    <td>offsets</td>
    <td>octet</td>
    <td colspan="8">0</td>
    <td colspan="8">1</td>
  </tr>

  <tr align="center">
		<td>octet</td>
		<td>bit</td>
		<td>7</td>
		<td>6</td>
		<td>5</td>
		<td>4</td>
		<td>3</td>
		<td>2</td>
		<td>1</td>
		<td>0</td>
		<td>7</td>
		<td>6</td>
		<td>5</td>
		<td>4</td>
		<td>3</td>
		<td>2</td>
		<td>1</td>
		<td>0</td>
  </tr>

  <tr align="center">
    <td>0</td>
    <td>0</td>
    <td colspan="16">payload length</td>
  </tr>

  <tr align="center">
    <td>2</td>
    <td>16</td>
    <td colspan="8">header length</td>
    <td colspan="4">source port</td>
    <td colspan="4">destination port</td>
  </tr>
</table>

This table describes header. Payload starts right after header.

**Payload Length (16 bits)**

Count of bytes that are related to actual payload.

**Header Length (8 bits)**

Count of bytes that are related to header. To offset to actual payload, you should use this value, not hard-coded fixed number based from table above.

**Source port (4 bits)**

Sender can specify for receiver where potential response is expected. Note that receiver still can use any value.

**Destination port (4 bits)**

Sender can specify for receiver, at receiver side handler at which location should handle sent data. Note that receiver still can use any handler.

## URL

STTP resources is a handlers. Handlers are identified and located on the network by URLs, using the URI scheme `sttp`.

Structure of STTP URL is `sttp://<IP>:<TCP PORT>/<HANDLER ID>`. Example: `sttp://192.168.1.235:4444/0`.
