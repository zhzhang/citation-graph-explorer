# Use python3 base image
FROM python:3

# Set working directory to /app
WORKDIR /app

# Copy source to container at /app
COPY . /app

# Install python dependencies
RUN pip install --trusted-host pypi.python.org -r requirements.txt

# Expose port 8000 for API
EXPOSE 8000

# Expore 3000 for front end
EXPOSE 3000

# Start API server
CMD ["gunicorn", "server:api", "--bind=0.0.0.0:8000"]
