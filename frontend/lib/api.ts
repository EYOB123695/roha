import axios from "axios" ;

import {useAuthStore} from  "@/store/useAuthStore"; 


const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "https://roha-production.up.railway.app";
export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true, // Sends cookies along with requests
});

// Request Interceptor: Automatically attach JWT token to Authorization header
api.interceptors.request.use(
  (config) => {
    // Read token directly from Zustand state if available
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);
// Response Interceptor: Catch 401 Unauthorized errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Optional: automatically logout on expired token
      const logout = useAuthStore.getState().logout;
      logout();
    }
    return Promise.reject(error);
  }
);

// Helper function to extract user-friendly error message from API response
export const getApiErrorMessage = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    return (
      error.response?.data?.error ||
      error.response?.data?.message ||
      error.message ||
      "An unexpected server error occurred."
    );
  }
  return "An unexpected error occurred.";
};








