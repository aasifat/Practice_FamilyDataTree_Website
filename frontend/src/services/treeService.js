import api from './api';

export const createFamilyTree = async (name) => {
  const response = await api.post('/trees', { name });
  return response.data;
};

export const getUserTrees = async () => {
  const response = await api.get('/trees');
  return response.data;
};

export const getFamilyTree = async (treeId) => {
  const response = await api.get(`/trees/${treeId}`);
  return response.data;
};

export const updateFamilyTree = async (treeId, data) => {
  const response = await api.put(`/trees/${treeId}`, data);
  return response.data;
};

export const deleteFamilyTree = async (treeId) => {
  const response = await api.delete(`/trees/${treeId}`);
  return response.data;
};

export const getTreeMembers = async (treeId) => {
  const response = await api.get(`/trees/${treeId}/members`);
  return response.data;
};

export const createPerson = async (treeId, personData) => {
  const response = await api.post(`/trees/${treeId}/members`, personData);
  return response.data;
};

export const updatePerson = async (personId, personData) => {
  const response = await api.put(`/members/${personId}`, personData);
  return response.data;
};

export const deletePerson = async (personId) => {
  const response = await api.delete(`/members/${personId}`);
  return response.data;
};

export const getPerson = async (personId) => {
  const response = await api.get(`/members/${personId}`);
  return response.data;
};

export const getChildren = async (personId) => {
  const response = await api.get(`/members/${personId}/children`);
  return response.data;
};

export const searchPeople = async (treeId, searchTerm) => {
  const response = await api.get(`/trees/${treeId}/members/search`, {
    params: { q: searchTerm },
  });
  return response.data;
};
