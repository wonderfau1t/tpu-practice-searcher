import axios from 'axios';
import { tg } from "../lib/telegram.ts";

export const API_URL = '';
const axiosInstance = axios.create({
  baseURL: API_URL,
});

axiosInstance.interceptors.request.use(
  (config) => {
    const accessToken = localStorage.getItem('accessToken');
    if (accessToken) {
      config.headers.Authorization = `Bearer ${accessToken}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      const initData = tg.initData;
      // const initData = "user=%7B%22id%22%3A508884173%2C%22first_name%22%3A%22wonderfau1t%22%2C%22last_name%22%3A%22%22%2C%22username%22%3A%22wonderrfau1t%22%2C%22language_code%22%3A%22ru%22%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2FxsX_1dpg_SjQyUK-1YJhrDZpDcvysMqWxIHy-y8gOak.svg%22%7D&chat_instance=-2504392250637437013&chat_type=private&auth_date=1744403963&signature=vUubuTC0aSxhLO26jiLHCM_gGrqKehEW11AMB2F9wrcO78RO0k2cfd8Wod50G1U_v70Qq6pZrUn2h2dP340IAg&hash=ac26d4c1be1c47f704dac579eb708ebb65b4256ce1f539ef434f44fb3222d2b7";
      try {
        const response = await axiosInstance.get('/auth', {
          headers: {
            Authorization: `tma ${initData}`,
          },
        });
        localStorage.setItem('accessToken', response.data.result.accessToken);
        originalRequest.headers.Authorization = `Bearer ${response.data.result.accessToken}`;
        return axiosInstance(originalRequest);
      } catch (error) {
        localStorage.removeItem('accessToken');
        window.location.href = '/'
        return Promise.reject(error);
      }
    }
    return Promise.reject(error);
  }
);

export default axiosInstance;