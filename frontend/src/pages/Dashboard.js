import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import {
  getUserTrees,
  createFamilyTree,
  deleteFamilyTree,
} from "../services/treeService";
import { useAuth } from "../context/AuthContext";
import {
  MdAdd,
  MdDelete,
  MdVisibility,
  MdDateRange,
  MdWarning,
} from "react-icons/md";
import "../styles/Dashboard.css";

const Dashboard = () => {
  const { user, loading: authLoading } = useAuth();
  const [trees, setTrees] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [newTreeName, setNewTreeName] = useState("");
  const [showForm, setShowForm] = useState(false);

  useEffect(() => {
    if (user) {
      loadTrees();
    }
  }, [user]);

  const loadTrees = async () => {
    try {
      setLoading(true);
      setError("");
      const data = await getUserTrees();
      setTrees(Array.isArray(data) ? data : []);
    } catch (err) {
      setError("Failed to load trees");
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTree = async (e) => {
    e.preventDefault();
    if (!newTreeName.trim()) {
      setError("Tree name is required");
      return;
    }

    try {
      setError("");
      await createFamilyTree(newTreeName);
      setNewTreeName("");
      setShowForm(false);
      await loadTrees();
    } catch (err) {
      setError("Failed to create tree");
    }
  };

  const handleDeleteTree = async (treeId) => {
    if (window.confirm("Are you sure? This action cannot be undone.")) {
      try {
        setError("");
        await deleteFamilyTree(treeId);
        await loadTrees();
      } catch (err) {
        setError("Failed to delete tree");
      }
    }
  };

  if (authLoading) {
    return (
      <div className="page-container">
        <div style={{ textAlign: "center", padding: "4rem 0" }}>
          <div className="loader"></div>
          <p style={{ marginTop: "1rem", color: "#6b7280" }}>Loading...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="page-container">
        <div style={{ maxWidth: "1200px", margin: "0 auto" }}>
          <div
            className="card"
            style={{ textAlign: "center", padding: "3rem" }}
          >
            <h2>Welcome to Family Tree Manager</h2>
            <p
              style={{
                marginTop: "1rem",
                fontSize: "1.1rem",
                color: "#6b7280",
              }}
            >
              Please{" "}
              <Link to="/login" style={{ fontWeight: 600 }}>
                login
              </Link>{" "}
              or{" "}
              <Link to="/signup" style={{ fontWeight: 600 }}>
                sign up
              </Link>{" "}
              to continue.
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="page-container">
      <div className="dashboard-container">
        {/* Header */}
        <div className="dashboard-header">
          <div>
            <h1>
              👋 Welcome,{" "}
              <span className="gradient-text">{user.first_name}</span>
            </h1>
            <p className="dashboard-subtitle">
              Explore and manage your family trees
            </p>
          </div>
          {!showForm && (
            <button
              onClick={() => setShowForm(true)}
              className="btn btn-primary"
            >
              <MdAdd size={18} /> Create Tree
            </button>
          )}
        </div>

        {/* Error Alert */}
        {error && (
          <div className="alert alert-error" style={{ marginBottom: "2rem" }}>
            <MdWarning size={20} />
            <span>{error}</span>
          </div>
        )}

        {/* Create Form */}
        {showForm && (
          <div className="create-form-card">
            <h3>Create a New Family Tree</h3>
            <form onSubmit={handleCreateTree} className="create-form">
              <div className="form-group">
                <input
                  type="text"
                  value={newTreeName}
                  onChange={(e) => setNewTreeName(e.target.value)}
                  placeholder="Enter tree name (e.g., Smith Family)"
                  className="form-input"
                  maxLength={50}
                  autoFocus
                />
                <p className="form-helper">
                  {newTreeName.length}/50 characters
                </p>
              </div>

              <div className="form-actions">
                <button
                  type="submit"
                  disabled={loading}
                  className="btn btn-primary"
                >
                  {loading ? (
                    <>
                      <span className="loader"></span>
                      Creating...
                    </>
                  ) : (
                    <>
                      <MdAdd size={18} /> Create
                    </>
                  )}
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setShowForm(false);
                    setNewTreeName("");
                  }}
                  className="btn btn-outline"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}

        {/* Trees Section */}
        <div className="trees-section">
          <h2>Your Family Trees</h2>

          {loading && !showForm ? (
            <div style={{ textAlign: "center", padding: "3rem 0" }}>
              <div className="loader"></div>
              <p style={{ marginTop: "1rem", color: "#6b7280" }}>
                Loading your trees...
              </p>
            </div>
          ) : trees.length === 0 ? (
            <div className="empty-state">
              <div className="empty-icon">🌳</div>
              <h3>No Family Trees Yet</h3>
              <p>Get started by creating your first family tree</p>
              {!showForm && (
                <button
                  onClick={() => setShowForm(true)}
                  className="btn btn-primary"
                  style={{ marginTop: "1.5rem" }}
                >
                  <MdAdd size={18} /> Create Your First Tree
                </button>
              )}
            </div>
          ) : (
            <div className="trees-grid">
              {trees.map((tree) => (
                <div key={tree.id} className="tree-card">
                  <div className="tree-card-header">
                    <h3>{tree.name}</h3>
                    <button
                      className="tree-delete-btn"
                      onClick={() => handleDeleteTree(tree.id)}
                      title="Delete tree"
                    >
                      <MdDelete size={18} />
                    </button>
                  </div>

                  <div className="tree-card-info">
                    <div className="info-item">
                      <MdDateRange size={16} />
                      <span>
                        {new Date(tree.created_at).toLocaleDateString("en-US", {
                          year: "numeric",
                          month: "short",
                          day: "numeric",
                        })}
                      </span>
                    </div>
                  </div>

                  <div className="tree-card-footer">
                    <Link
                      to={`/tree/${tree.id}`}
                      className="btn btn-primary"
                      style={{ flex: 1 }}
                    >
                      <MdVisibility size={16} /> View Tree
                    </Link>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
