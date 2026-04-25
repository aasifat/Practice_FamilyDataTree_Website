import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import {
  MdLogout,
  MdPerson,
  MdDashboard,
  MdLogin,
  MdMenu,
} from "react-icons/md";
import { MdClose } from "react-icons/md";
import "../styles/Navbar.css";

const Navbar = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const handleLogout = () => {
    logout();
    navigate("/login");
    setMobileMenuOpen(false);
  };

  return (
    <nav className="navbar">
      <div className="navbar-container">
        <Link to="/" className="navbar-logo">
          <span className="logo-icon">🌳</span>
          <span className="logo-text">Family Tree</span>
        </Link>

        <button
          className="mobile-toggle"
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
        >
          {mobileMenuOpen ? <MdClose size={24} /> : <MdMenu size={24} />}
        </button>

        <ul className={`navbar-menu ${mobileMenuOpen ? "active" : ""}`}>
          <li className="navbar-item">
            <Link
              to="/"
              className="navbar-link"
              onClick={() => setMobileMenuOpen(false)}
            >
              <MdDashboard size={18} /> Dashboard
            </Link>
          </li>
          <li className="navbar-item">
            <Link
              to="/trees"
              className="navbar-link"
              onClick={() => setMobileMenuOpen(false)}
            >
              Trees
            </Link>
          </li>
          <li className="navbar-item">
            <Link
              to="/contact"
              className="navbar-link"
              onClick={() => setMobileMenuOpen(false)}
            >
              Contact
            </Link>
          </li>

          {user ? (
            <>
              <li className="navbar-item">
                <Link
                  to="/profile"
                  className="navbar-link"
                  onClick={() => setMobileMenuOpen(false)}
                >
                  <MdPerson size={18} /> Profile
                </Link>
              </li>
              <li className="navbar-item">
                <button className="navbar-logout-btn" onClick={handleLogout}>
                  <MdLogout size={18} /> Logout
                </button>
              </li>
              <li className="navbar-user-info">
                <span className="user-avatar">{user.first_name?.[0]}</span>
                <span className="user-name">{user.first_name}</span>
              </li>
            </>
          ) : (
            <>
              <li className="navbar-item">
                <Link
                  to="/login"
                  className="navbar-link"
                  onClick={() => setMobileMenuOpen(false)}
                >
                  <MdLogin size={18} /> Login
                </Link>
              </li>
              <li className="navbar-item">
                <Link
                  to="/signup"
                  className="navbar-btn navbar-signup-btn"
                  onClick={() => setMobileMenuOpen(false)}
                >
                  Sign Up
                </Link>
              </li>
            </>
          )}
        </ul>
      </div>
    </nav>
  );
};

export default Navbar;
