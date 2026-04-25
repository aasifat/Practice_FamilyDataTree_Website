import React, { useState } from "react";
import { Link } from "react-router-dom";
import { forgotPassword } from "../services/authService";
import { MdEmail, MdArrowBack } from "react-icons/md";
import "../styles/Auth.css";

const ForgotPassword = () => {
  const [email, setEmail] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setMessage("");
    setLoading(true);

    try {
      const result = await forgotPassword(email);
      setMessage(
        result.message || "If the email exists, a reset link has been sent.",
      );
      setEmail("");
    } catch (err) {
      setError(err.response?.data?.error || "Failed to request password reset");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-background">
        <div className="auth-gradient-1"></div>
        <div className="auth-gradient-2"></div>
      </div>

      <div className="auth-content">
        <div className="auth-card">
          <div className="auth-header">
            <h1>Forgot Password?</h1>
            <p>We'll send you a link to reset it</p>
          </div>

          {message && (
            <div className="alert alert-success">
              <span>{message}</span>
            </div>
          )}

          {error && (
            <div className="alert alert-error">
              <span>{error}</span>
            </div>
          )}

          <form onSubmit={handleSubmit} className="auth-form">
            <div className="form-group">
              <label className="form-label">Email Address</label>
              <div className="input-wrapper">
                <MdEmail className="input-icon" size={18} />
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="you@example.com"
                  required
                  className="form-input"
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="btn btn-primary btn-block"
            >
              {loading ? (
                <>
                  <span className="loader"></span>
                  Sending Link...
                </>
              ) : (
                "Send Reset Link"
              )}
            </button>
          </form>

          <div className="auth-divider">
            <span>Remember your password?</span>
          </div>

          <Link to="/login" className="btn btn-outline btn-block">
            <MdArrowBack size={16} /> Back to Login
          </Link>
        </div>

        <div className="auth-info">
          <div className="info-card">
            <div className="info-icon">🔐</div>
            <h3>Secure Reset</h3>
            <p>Your account security is our top priority</p>
          </div>
          <div className="info-card">
            <div className="info-icon">⚡</div>
            <h3>Quick Process</h3>
            <p>Reset your password in just a few seconds</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ForgotPassword;
