import { describe, expect, test } from "bun:test";

import * as _ from "./_";

describe("enum smoke", () => {
    test("AccountStatus validator accepts generated values", () => {
        expect(_.IsAccountStatus(_.AccountStatus.Offline)).toBe(true);
        expect(_.IsAccountStatus(_.AccountStatus.Deleted)).toBe(true);
        expect(_.IsAccountStatus(255 as any)).toBe(false);
    });
    test("Type validator accepts generated values", () => {
        expect(_.IsType(_.Type.Sim)).toBe(true);
        expect(_.IsType(_.Type.Recharge)).toBe(true);
        expect(_.IsType(255 as any)).toBe(false);
    });
    test("Status validator accepts generated values", () => {
        expect(_.IsStatus(_.Status.Ok)).toBe(true);
        expect(_.IsStatus(_.Status.One)).toBe(true);
        expect(_.IsStatus(255 as any)).toBe(false);
    });
    test("StatusA validator accepts generated values", () => {
        expect(_.IsStatusA(_.StatusA.Ok)).toBe(true);
        expect(_.IsStatusA(_.StatusA.Seven)).toBe(true);
        expect(_.IsStatusA(255 as any)).toBe(false);
    });
    test("ItemStatus validator accepts generated values", () => {
        expect(_.IsItemStatus(_.ItemStatus.Offline)).toBe(true);
        expect(_.IsItemStatus(_.ItemStatus.Online)).toBe(true);
        expect(_.IsItemStatus(255 as any)).toBe(false);
    });
    test("SimPickPhone validator accepts generated values", () => {
        expect(_.IsSimPickPhone(_.SimPickPhone.No)).toBe(true);
        expect(_.IsSimPickPhone(_.SimPickPhone.Abcc)).toBe(true);
        expect(_.IsSimPickPhone(255 as any)).toBe(false);
    });
    test("SimOperator validator accepts generated values", () => {
        expect(_.IsSimOperator(_.SimOperator.Zz)).toBe(true);
        expect(_.IsSimOperator(_.SimOperator.B)).toBe(true);
        expect(_.IsSimOperator(255 as any)).toBe(false);
    });
    test("OrderStatus validator accepts generated values", () => {
        expect(_.IsOrderStatus(_.OrderStatus.Pending)).toBe(true);
        expect(_.IsOrderStatus(_.OrderStatus.Settled)).toBe(true);
        expect(_.IsOrderStatus(255 as any)).toBe(false);
    });
});
