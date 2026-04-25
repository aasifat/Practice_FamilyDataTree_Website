import React, { useState } from "react";
import { useAuth } from "../context/AuthContext";
import {
  MdEmail,
  MdDateRange,
  MdSecurity,
  MdEdit,
  MdCheckCircle,
} from "react-icons/md";
import "../styles/Profile.css";

const Profile = () => {
  const { user } = useAuth();
  const [isEditing, setIsEditing] = useState(false);

  if (!user) {
    return (
      <div className="page-container">
        <div className="profile-container">
          <div
            className="card"
            style={{ textAlign: "center", padding: "3rem" }}
          >
            <p style={{ fontSize: "1.1rem", color: "#6b7280" }}>
              Please login to view your profile
            </p>
          </div>
        </div>
      </div>
    );
  }

  const memberSince = new Date(user.created_at).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  return (
    <div className="page-container">
      <div className="profile-container">
        {/* Header */}
        <div className="profile-header">
          <h1>Your Profile</h1>
          <p className="profile-subtitle">Manage your account information</p>
        </div>

        {/* Main Card */}
        <div className="profile-card">
          {/* Profile Info Section */}
          <div className="profile-info-section">
            {/* Avatar */}
            <div className="profile-avatar">
              <span>{user.first_name?.[0]}</span>
            </div>

            {/* User Info */}
            <div className="profile-details">
              <h2>{`${user.first_name} ${user.last_name}`}</h2>
              <p className="profile-role">Member</p>

              {/* Info Grid */}
              <div className="info-grid">
                <div className="info-card">
                  <div className="info-header">
                    <MdEmail className="info-icon" size={20} />
                    <span className="info-label">Email</span>
                  </div>
                  <p className="info-value">{user.email}</p>
                </div>

                <div className="info-card">
                  <div className="info-header">
                    <MdSecurity className="info-icon" size={20} />
                    <span className="info-label">Role</span>
                  </div>
                  <p className="info-value">{user.role}</p>
                </div>

                <div className="info-card">
                  <div className="info-header">
                    <MdDateRange className="info-icon" size={20} />
                    <span className="info-label">Member Since</span>
                  </div>
                  <p className="info-value">{memberSince}</p>
                </div>

                <div className="info-card">
                  <div className="info-header">
                    <MdCheckCircle className="info-icon" size={20} />
                    <span className="info-label">Status</span>
                  </div>
                  <p className="info-value">
                    <span className="status-badge active">Active</span>
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Action Section */}
          <div className="profile-actions">
            <button
              className="btn btn-primary"
              onClick={() => setIsEditing(!isEditing)}
            >
              <MdEdit size={18} /> {isEditing ? "Cancel" : "Edit Profile"}
            </button>
          </div>
        </div>

        {/* Additional Info Cards */}
        <div className="additional-cards">
          <div className="feature-card">
            <div className="feature-icon">🔐</div>
            <h3>Security</h3>
            <p>Secure your account with a strong password</p>
            <a href="/forgot-password" className="feature-link">
              Change Password →
            </a>
          </div>

          <div className="feature-card">
            <div className="feature-icon">🌳</div>
            <h3>Family Trees</h3>
            <p>View and manage all your family trees</p>
            <a href="/" className="feature-link">
              Go to Trees →
            </a>
          </div>

          <div className="feature-card">
            <div className="feature-icon">📧</div>
            <h3>Notifications</h3>
            <p>Stay updated with email notifications</p>
            <button className="feature-link">Manage Settings →</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Profile;
