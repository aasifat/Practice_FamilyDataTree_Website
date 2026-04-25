import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import {
  getFamilyTree,
  createPerson,
  deletePerson,
  searchPeople,
} from "../services/treeService";
import TreeVisualization from "../components/TreeVisualization";

const TreeDetail = () => {
  const { id } = useParams();
  const [tree, setTree] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showForm, setShowForm] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [searchResults, setSearchResults] = useState([]);
  const [formData, setFormData] = useState({
    name: "",
    gender: "male",
    fatherId: "",
    motherId: "",
    spouseId: "",
  });

  const loadTree = React.useCallback(async () => {
    try {
      setLoading(true);
      const data = await getFamilyTree(id);
      setTree(data);
    } catch (err) {
      setError("Failed to load tree");
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    loadTree();
  }, [id, loadTree]);

  const handleAddMember = async (e) => {
    e.preventDefault();
    try {
      await createPerson(id, {
        name: formData.name,
        gender: formData.gender,
        father_id: formData.fatherId ? parseInt(formData.fatherId) : null,
        mother_id: formData.motherId ? parseInt(formData.motherId) : null,
        spouse_id: formData.spouseId ? parseInt(formData.spouseId) : null,
      });
      setFormData({
        name: "",
        gender: "male",
        fatherId: "",
        motherId: "",
        spouseId: "",
      });
      setShowForm(false);
      await loadTree();
    } catch (err) {
      setError(
        "Failed to add member: " +
          (err.message ||
            "Please ensure both father and mother are specified for children"),
      );
    }
  };

  const handleDeleteMember = async (personId) => {
    if (window.confirm("Are you sure?")) {
      try {
        await deletePerson(personId);
        await loadTree();
      } catch (err) {
        setError("Failed to delete member");
      }
    }
  };

  const handleSearch = async (e) => {
    if (e.key === "Enter" && searchTerm.trim()) {
      try {
        const results = await searchPeople(id, searchTerm);
        setSearchResults(results);
      } catch (err) {
        setError("Search failed");
      }
    }
  };

  if (loading) return <div style={styles.container}>Loading...</div>;
  if (!tree) return <div style={styles.container}>Tree not found</div>;

  return (
    <div style={styles.container}>
      <div style={styles.content}>
        <h1>{tree.name}</h1>

        {error && <div style={styles.error}>{error}</div>}

        <div style={styles.toolbar}>
          <input
            type="text"
            placeholder="Search members..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onKeyPress={handleSearch}
            style={styles.searchInput}
          />
          <button onClick={() => setShowForm(!showForm)} style={styles.button}>
            {showForm ? "Cancel" : "+ Add Member"}
          </button>
        </div>

        {showForm && (
          <form onSubmit={handleAddMember} style={styles.form}>
            <input
              type="text"
              placeholder="Member name"
              value={formData.name}
              onChange={(e) =>
                setFormData({ ...formData, name: e.target.value })
              }
              required
              style={styles.input}
            />
            <select
              value={formData.gender}
              onChange={(e) =>
                setFormData({ ...formData, gender: e.target.value })
              }
              style={styles.input}
            >
              <option value="male">Male</option>
              <option value="female">Female</option>
            </select>

            <select
              value={formData.fatherId}
              onChange={(e) =>
                setFormData({ ...formData, fatherId: e.target.value })
              }
              style={styles.input}
              title="Select father (if adding a child)"
            >
              <option value="">No Father (Root Level)</option>
              {tree.members &&
                tree.members
                  .filter((m) => m.gender === "male")
                  .map((member) => (
                    <option key={member.id} value={member.id}>
                      {member.name}
                    </option>
                  ))}
            </select>

            <select
              value={formData.motherId}
              onChange={(e) =>
                setFormData({ ...formData, motherId: e.target.value })
              }
              style={styles.input}
              title="Select mother (if adding a child)"
            >
              <option value="">No Mother (Root Level)</option>
              {tree.members &&
                tree.members
                  .filter((m) => m.gender === "female")
                  .map((member) => (
                    <option key={member.id} value={member.id}>
                      {member.name}
                    </option>
                  ))}
            </select>

            <select
              value={formData.spouseId}
              onChange={(e) =>
                setFormData({ ...formData, spouseId: e.target.value })
              }
              style={styles.input}
              title="Optional: Select spouse"
            >
              <option value="">No Spouse</option>
              {tree.members &&
                tree.members.map((member) => (
                  <option key={member.id} value={member.id}>
                    {member.name}
                  </option>
                ))}
            </select>

            <button type="submit" style={styles.submitBtn}>
              Add Member
            </button>
          </form>
        )}

        {searchResults.length > 0 && (
          <div style={styles.searchResultsContainer}>
            <h3>Search Results:</h3>
            <div style={styles.membersList}>
              {searchResults.map((member) => (
                <div key={member.id} style={styles.memberCard}>
                  <p>
                    <strong>{member.name}</strong> ({member.gender})
                  </p>
                </div>
              ))}
            </div>
          </div>
        )}

        <TreeVisualization tree={tree} onDeleteMember={handleDeleteMember} />

        {tree.members && tree.members.length > 0 && (
          <div style={styles.membersSection}>
            <h2>All Members</h2>
            <div style={styles.membersList}>
              {tree.members.map((member) => (
                <div key={member.id} style={styles.memberCard}>
                  <p>
                    <strong>{member.name}</strong>
                  </p>
                  <p
                    style={{
                      fontSize: "0.85rem",
                      color: "#666",
                      margin: "0.3rem 0",
                    }}
                  >
                    Gender: {member.gender}
                  </p>
                  {member.father_id && (
                    <p
                      style={{
                        fontSize: "0.85rem",
                        color: "#666",
                        margin: "0.3rem 0",
                      }}
                    >
                      Father ID: {member.father_id}
                    </p>
                  )}
                  {member.mother_id && (
                    <p
                      style={{
                        fontSize: "0.85rem",
                        color: "#666",
                        margin: "0.3rem 0",
                      }}
                    >
                      Mother ID: {member.mother_id}
                    </p>
                  )}
                  {member.spouse_id && (
                    <p
                      style={{
                        fontSize: "0.85rem",
                        color: "#e91e63",
                        margin: "0.3rem 0",
                      }}
                    >
                      ♥ Spouse ID: {member.spouse_id}
                    </p>
                  )}
                  <div style={styles.actions}>
                    <a href={`/member/${member.id}`} style={styles.link}>
                      View
                    </a>
                    <button
                      onClick={() => handleDeleteMember(member.id)}
                      style={styles.deleteBtn}
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

const styles = {
  container: {
    minHeight: "calc(100vh - 80px)",
    backgroundColor: "#f5f5f5",
    padding: "2rem",
  },
  content: {
    maxWidth: "1200px",
    margin: "0 auto",
    backgroundColor: "white",
    padding: "2rem",
    borderRadius: "8px",
  },
  toolbar: {
    display: "flex",
    gap: "1rem",
    marginBottom: "1rem",
  },
  searchInput: {
    flex: 1,
    padding: "0.75rem",
    border: "1px solid #ddd",
    borderRadius: "4px",
  },
  button: {
    padding: "0.75rem 1.5rem",
    backgroundColor: "#27ae60",
    color: "white",
    border: "none",
    borderRadius: "4px",
    cursor: "pointer",
  },
  form: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fit, minmax(150px, 1fr))",
    gap: "1rem",
    marginBottom: "2rem",
    padding: "1.5rem",
    backgroundColor: "#ecf0f1",
    borderRadius: "4px",
    border: "2px solid #bdc3c7",
  },
  input: {
    padding: "0.75rem",
    border: "1px solid #ddd",
    borderRadius: "4px",
  },
  submitBtn: {
    padding: "0.75rem",
    backgroundColor: "#27ae60",
    color: "white",
    border: "none",
    borderRadius: "4px",
    cursor: "pointer",
  },
  error: {
    color: "#e74c3c",
    padding: "1rem",
    backgroundColor: "#fadbd8",
    borderRadius: "4px",
    marginBottom: "1rem",
  },
  searchResultsContainer: {
    marginBottom: "2rem",
  },
  membersSection: {
    marginTop: "2rem",
  },
  membersList: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fill, minmax(250px, 1fr))",
    gap: "1rem",
    marginTop: "1rem",
  },
  memberCard: {
    backgroundColor: "#ecf0f1",
    padding: "1rem",
    borderRadius: "4px",
  },
  actions: {
    display: "flex",
    gap: "0.5rem",
    marginTop: "1rem",
  },
  link: {
    padding: "0.5rem 1rem",
    backgroundColor: "#3498db",
    color: "white",
    textDecoration: "none",
    borderRadius: "4px",
    cursor: "pointer",
  },
  deleteBtn: {
    padding: "0.5rem 1rem",
    backgroundColor: "#e74c3c",
    color: "white",
    border: "none",
    borderRadius: "4px",
    cursor: "pointer",
  },
};

export default TreeDetail;
