import React, { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { resetPassword } from "../services/authService";
import {
  MdLock,
  MdVisibility,
  MdVisibilityOff,
  MdCheckCircle,
} from "react-icons/md";
import "../styles/Auth.css";

const ResetPassword = () => {
  const [searchParams] = useSearchParams();
  const [token, setToken] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    const tokenValue = searchParams.get("token");
    if (tokenValue) {
      setToken(tokenValue);
    }
  }, [searchParams]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setMessage("");

    if (!token) {
      setError("Reset token is missing.");
      return;
    }

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    if (password.length < 6) {
      setError("Password must be at least 6 characters");
      return;
    }

    setLoading(true);

    try {
      const result = await resetPassword(token, password);
      setMessage(result.message || "Your password was updated successfully.");
      setPassword("");
      setConfirmPassword("");
      setSuccess(true);
    } catch (err) {
      setError(err.response?.data?.error || "Failed to reset password");
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="auth-container">
        <div className="auth-background">
          <div className="auth-gradient-1"></div>
          <div className="auth-gradient-2"></div>
        </div>

        <div className="auth-content">
          <div className="auth-card" style={{ textAlign: "center" }}>
            <div style={{ fontSize: "4rem", marginBottom: "1rem" }}>
              <MdCheckCircle size={64} color="#10b981" />
            </div>
            <h1 style={{ color: "#10b981", marginBottom: "1rem" }}>
              Password Reset!
            </h1>
            <p
              style={{
                marginBottom: "2rem",
                fontSize: "1.05rem",
                color: "#6b7280",
              }}
            >
              Your password has been successfully reset. You can now login with
              your new password.
            </p>
            <Link to="/login" className="btn btn-primary btn-block">
              Go to Login
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="auth-container">
      <div className="auth-background">
        <div className="auth-gradient-1"></div>
        <div className="auth-gradient-2"></div>
      </div>

      <div className="auth-content">
        <div className="auth-card">
          <div className="auth-header">
            <h1>Create New Password</h1>
            <p>Choose a strong password to secure your account</p>
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
              <label className="form-label">New Password</label>
              <div className="input-wrapper">
                <MdLock className="input-icon" size={18} />
                <input
                  type={showPassword ? "text" : "password"}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••"
                  required
                  minLength={6}
                  className="form-input"
                />
                <button
                  type="button"
                  className="password-toggle"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? (
                    <MdVisibilityOff size={18} />
                  ) : (
                    <MdVisibility size={18} />
                  )}
                </button>
              </div>
            </div>

            <div className="form-group">
              <label className="form-label">Confirm Password</label>
              <div className="input-wrapper">
                <MdLock className="input-icon" size={18} />
                <input
                  type={showConfirmPassword ? "text" : "password"}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  placeholder="••••••••"
                  required
                  minLength={6}
                  className="form-input"
                />
                <button
                  type="button"
                  className="password-toggle"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                >
                  {showConfirmPassword ? (
                    <MdVisibilityOff size={18} />
                  ) : (
                    <MdVisibility size={18} />
                  )}
                </button>
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
                  Resetting...
                </>
              ) : (
                "Reset Password"
              )}
            </button>
          </form>

          <div className="auth-divider">
            <span>Ready to login?</span>
          </div>

          <Link to="/login" className="btn btn-outline btn-block">
            Go to Login
          </Link>
        </div>

        <div className="auth-info">
          <div className="info-card">
            <div className="info-icon">🛡️</div>
            <h3>Strong Security</h3>
            <p>Protect your account with a strong password</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ResetPassword;
