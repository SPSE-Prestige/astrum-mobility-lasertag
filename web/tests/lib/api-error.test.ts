import { describe, it, expect } from "vitest";
import { ApiError } from "@/lib/api/client";

describe("ApiError", () => {
  it("sets status, code, and data", () => {
    const err = new ApiError(404, "not_found", { detail: "missing" });
    expect(err.status).toBe(404);
    expect(err.code).toBe("not_found");
    expect(err.data).toEqual({ detail: "missing" });
    expect(err.message).toBe("not_found");
    expect(err.name).toBe("ApiError");
  });

  it("isAuthError is true for 401", () => {
    expect(new ApiError(401, "unauthorized").isAuthError).toBe(true);
    expect(new ApiError(403, "forbidden").isAuthError).toBe(false);
  });

  it("isNotFound is true for 404", () => {
    expect(new ApiError(404, "not_found").isNotFound).toBe(true);
    expect(new ApiError(400, "bad_request").isNotFound).toBe(false);
  });

  it("isConflict is true for 409", () => {
    expect(new ApiError(409, "conflict").isConflict).toBe(true);
  });

  it("isServerError is true for 5xx", () => {
    expect(new ApiError(500, "internal").isServerError).toBe(true);
    expect(new ApiError(502, "bad_gateway").isServerError).toBe(true);
    expect(new ApiError(499, "client").isServerError).toBe(false);
  });

  it("falls back to HTTP status when code is empty", () => {
    const err = new ApiError(500, "");
    expect(err.message).toBe("HTTP 500");
  });
});
