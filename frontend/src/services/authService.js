import api from "./api";
import { setAuthToken, removeAuthToken } from "../config";

export const signup = async (userData) => {
  const response = await api.post("/auth/signup", userData);
  if (response.data.token) {
    setAuthToken(response.data.token);
  }
  return response.data;
};

export const requestOTPForSignup = async (userData) => {
  const response = await api.post("/auth/request-otp-signup", userData);
  return response.data;
};

export const verifyOTPAndSignup = async (data) => {
  const response = await api.post("/auth/verify-otp-signup", data);
  if (response.data.token) {
    setAuthToken(response.data.token);
  }
  return response.data;
};

export const login = async (credentials) => {
  const response = await api.post("/auth/login", credentials);
  if (response.data.token) {
    setAuthToken(response.data.token);
  }
  return response.data;
};

export const logout = () => {
  removeAuthToken();
};

export const getProfile = async () => {
  const response = await api.get("/auth/profile");
  return response.data;
};

export const forgotPassword = async (email) => {
  const response = await api.post("/auth/forgot-password", { email });
  return response.data;
};

export const resetPassword = async (token, password) => {
  const response = await api.post("/auth/reset-password", { token, password });
  return response.data;
};
