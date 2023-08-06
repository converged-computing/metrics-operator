# The Metrics Operator

Welcome to the Metrics Operator Documentation!

The Metrics Operator is a Kubernetes Cluster [Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
that you can install to your cluster to easily measure different aspects of performance, including (but not limited to):

 - I/O metrics to assess storage options
 - system performance metrics (memory, cpu, etc.)
 - network or other custom metrics / timings

For this project, we aim to provide the following:

1. A catalog of pre-defined metrics and associated containers you can quickly use
2. A number of containerized HPC applications and experiments to demonstrate using the operator.

We wanted to create this operator because we didn't have a solid understanding
of our application performance (typically from the high performance computing space)
on Kubernetes. We also didn't have a good catalog of sample applications, and
wanted to provide this to the community.

The Metrics Operator is currently 🚧️ Under Construction! 🚧️
This is a *converged computing* project that aims
to unite the worlds and technologies typical of cloud computing and
high performance computing.

To get started, check out the links below!
Would you like to request a feature or contribute?
[Open an issue](https://github.com/converged-computing/metrics-operator/issues).

```{toctree}
:caption: Getting Started
:maxdepth: 2
getting_started/index.md
development/index.md
```

```{toctree}
:caption: About
:maxdepth: 2
about/index.md
```

<script>
// This is a small hack to populate empty sidebar with an image!
document.addEventListener('DOMContentLoaded', function () {
    var currentNode = document.querySelector('.md-sidebar__scrollwrap');
    currentNode.outerHTML =
	'<div class="md-sidebar__scrollwrap">' +
		'<img style="width:100%" src="_static/images/the-metrics-operator.png"/>' +

	'</div>';
}, false);

</script>
