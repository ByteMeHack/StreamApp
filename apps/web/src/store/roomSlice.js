import { createAsyncThunk, createSlice, isAnyOf } from "@reduxjs/toolkit";
import { createRoom, getUserRooms } from "../api";

export const setRooms = createAsyncThunk("rooms/setRooms", async (userId) => {
  const rooms = await getUserRooms(userId);
  return rooms;
});

export const addRoom = createAsyncThunk("rooms/addRoom", async (room) => {
  try {
    const rooms = await createRoom(
      room.name,
      room.private,
      room.password,
      room.owner_id
    );
    return rooms;
  } catch (err) {
    throw err;
  }
});

const roomSlice = createSlice({
  name: "roomSlice",
  initialState: {
    rooms: [],
    loading: false,
  },
  reducers: {},

  extraReducers: (builder) => {
    builder.addCase(setRooms.fulfilled, (state, action) => {
      state.loading = false;
      state.rooms = action.payload;
    });
    builder.addCase(addRoom.fulfilled, (state, action) => {
      state.loading = false;
      state.rooms.push(action.payload);
    });
    builder.addMatcher(
      isAnyOf(setRooms.pending, addRoom.pending),
      (state, action) => {
        state.loading = true;
      }
    );
    builder.addMatcher(
      isAnyOf(setRooms.rejected, addRoom.rejected),
      (state, action) => {
        state.loading = false;
      }
    );
  },
});

export default roomSlice.reducer;
