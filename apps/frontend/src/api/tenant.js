import axios from 'axios';

// Create an Axios instance for tenant API
const tenantApi = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:3000/api/v1',
  timeout: 5000,
});

export const fetchTenantConfig = async (hostname) => {
  try {
    // We send the hostname to the backend so it knows which tenant's config to return
    // In local dev, it might be localhost, so backend should have a fallback
    const response = await tenantApi.get(`/tenants/config`, {
      params: { hostname }
    });
    return response.data;
  } catch (error) {
    // If backend isn't ready or fails, return null to fallback to default theme
    console.warn("Could not fetch tenant config, using defaults.");
    return null;
  }
};

export const updateTenantConfig = async (config) => {
  try {
    const response = await tenantApi.put(`/tenants/config`, config, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}` // assuming auth token
      }
    });
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.message || 'Failed to update theme config');
  }
};
