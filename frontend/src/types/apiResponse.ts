export interface ApiResponse<T> {
  status: 'OK' | 'ERROR';
  result: T;
}