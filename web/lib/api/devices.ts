import { httpClient } from "./client";
import type { DeviceResponse } from "./types";

const enc = encodeURIComponent;

export const deviceApi = {
  list: () => httpClient.get<DeviceResponse[]>("/devices"),

  listAvailable: () => httpClient.get<DeviceResponse[]>("/devices/available"),
};
