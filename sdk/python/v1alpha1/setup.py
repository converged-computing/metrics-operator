"""
    metrics operator

    Python SDK for submitting CRD for metrics and parsing.
    This does not attempt to automate the API generation -
    we leave that up to the user (and provide examples)
    but rather provides classes for parsing data outputs.
"""

import os

from setuptools import find_packages, setup  # noqa: H301

# Make sure everything is relative to setup.py
install_path = os.path.dirname(os.path.abspath(__file__))
os.chdir(install_path)

DESCRIPTION = "Python helpers for the Metrics Operator"
# Try to read description, otherwise fallback to short description
try:
    with open(os.path.join("README.md")) as filey:
        LONG_DESCRIPTION = filey.read()
except Exception:
    LONG_DESCRIPTION = DESCRIPTION

################################################################################
# MAIN #########################################################################
################################################################################

if __name__ == "__main__":
    setup(
        name="metricsoperator",
        version="0.0.19",
        author="Vanessasaurus",
        author_email="vsoch@users.noreply.github.com",
        maintainer="Vanessasaurus",
        packages=find_packages(),
        include_package_data=True,
        zip_safe=False,
        url="https://github.com/converged-computing/metrics-operator/tree/main/python-sdk/v1alpha1",
        license="MIT",
        description=DESCRIPTION,
        long_description=LONG_DESCRIPTION,
        long_description_content_type="text/markdown",
        keywords="metrics-operator,hpc,kubernetes,metrics,storage,applications",
        setup_requires=["pytest-runner"],
        install_requires=["kubernetes", "requests", "pyyaml"],
        tests_require=["pytest", "pytest-cov"],
        classifiers=[
            "Intended Audience :: Science/Research",
            "Intended Audience :: Developers",
            "License :: OSI Approved :: Apache Software License",
            "Programming Language :: C",
            "Programming Language :: Python",
            "Topic :: Software Development",
            "Topic :: Scientific/Engineering",
            "Operating System :: Unix",
            "Programming Language :: Python :: 3.7",
        ],
    )
