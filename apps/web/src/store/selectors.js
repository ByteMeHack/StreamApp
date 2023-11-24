import { createSelector } from "@reduxjs/toolkit";

export const userSelector = createSelector(
  (state) => state.userSlice.user,
  (user) => user
);

export const roomsSelector = createSelector(
  (state) => state.roomSlice.rooms,
  (rooms) => rooms
);
